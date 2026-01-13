package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/ultima-linea/scraper/internal/models"
)

const (
	baseURL      = "https://huang-66.x.yupoo.com"
	galleryURL   = baseURL + "/albums?tab=gallery"
	maxPages     = 45 // Número total de páginas según el análisis
	imagesPerAlbum = 3  // Solo las primeras 3 imágenes
)

type YupooScraper struct {
	collector *colly.Collector
	albums    []models.YupooAlbum
	stats     models.ScraperStats
}

// NewYupooScraper crea una nueva instancia del scraper
func NewYupooScraper() *YupooScraper {
	c := colly.NewCollector(
		colly.AllowedDomains("huang-66.x.yupoo.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	// Rate limiting para no saturar el servidor
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 5,        // Aumentado de 2 a 5 requests paralelos
		Delay:       500 * time.Millisecond, // Reducido de 2s a 0.5s
	})

	return &YupooScraper{
		collector: c,
		albums:    make([]models.YupooAlbum, 0),
		stats: models.ScraperStats{
			StartTime: time.Now(),
		},
	}
}

// ScrapeAllPages recorre todas las páginas de la galería
func (s *YupooScraper) ScrapeAllPages(startPage, endPage int) error {
	log.Printf("Iniciando scraping de páginas %d a %d...\n", startPage, endPage)

	// Validar rango de páginas
	if startPage < 1 || endPage > maxPages || startPage > endPage {
		return fmt.Errorf("rango de páginas inválido: %d-%d (máximo: %d)", startPage, endPage, maxPages)
	}

	s.stats.TotalPages = endPage - startPage + 1

	// Configurar el scraper para el listado de álbumes
	s.setupAlbumListingCollector()

	// Recorrer cada página
	for page := startPage; page <= endPage; page++ {
		url := fmt.Sprintf("%s&page=%d", galleryURL, page)
		log.Printf("Scraping página %d/%d: %s\n", page, endPage, url)

		err := s.collector.Visit(url)
		if err != nil {
			log.Printf("Error al visitar página %d: %v\n", page, err)
			s.stats.FailedScans++
		}
	}

	s.collector.Wait()

	s.stats.EndTime = time.Now()
	s.stats.Duration = s.stats.EndTime.Sub(s.stats.StartTime).String()
	s.stats.TotalAlbums = len(s.albums)

	log.Printf("Scraping completado: %d álbumes encontrados\n", s.stats.TotalAlbums)

	return nil
}

// setupAlbumListingCollector configura el collector para el listado de álbumes
func (s *YupooScraper) setupAlbumListingCollector() {
	// Buscar cada enlace de álbum en la galería
	// Los álbumes son enlaces con formato /albums/[ID]?uid=1
	s.collector.OnHTML("a[href*='/albums/']", func(e *colly.HTMLElement) {
		album := models.YupooAlbum{
			ScrapedAt: time.Now(),
		}

		// Extraer el enlace del álbum
		link := e.Attr("href")
		if link == "" {
			return
		}

		// El href es relativo, construir URL completa
		album.AlbumURL = baseURL + link

		// Extraer ID del álbum desde la URL
		re := regexp.MustCompile(`/albums/(\d+)`)
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			album.ID = matches[1]
		} else {
			return // Si no tiene ID, no es un álbum válido
		}

		// Extraer título del atributo title del enlace
		album.Title = strings.TrimSpace(e.Attr("title"))

		// Si no hay title, usar el texto del enlace
		if album.Title == "" {
			album.Title = strings.TrimSpace(e.Text)
		}

		// Extraer el número de imágenes del texto del enlace
		// El formato es: "número\ntítulo" o solo el número
		text := strings.TrimSpace(e.Text)
		lines := strings.Split(text, "\n")
		if len(lines) > 0 {
			firstLine := strings.TrimSpace(lines[0])
			count, err := strconv.Atoi(firstLine)
			if err == nil {
				album.ImageCount = count
			}
		}

		// Extraer la página actual desde la URL
		currentURL := e.Request.URL.String()
		re = regexp.MustCompile(`page=(\d+)`)
		matches = re.FindStringSubmatch(currentURL)
		if len(matches) > 1 {
			pageNum, _ := strconv.Atoi(matches[1])
			album.PageNumber = pageNum
		} else {
			album.PageNumber = 1 // Primera página si no hay parámetro
		}

		// Solo agregar si tiene datos válidos
		if album.ID != "" && album.Title != "" {
			s.albums = append(s.albums, album)
			s.stats.SuccessfulScans++
			log.Printf("  ✓ Álbum encontrado: %s (ID: %s, %d imágenes)\n", album.Title, album.ID, album.ImageCount)
		}
	})

	s.collector.OnError(func(r *colly.Response, err error) {
		log.Printf("Error al hacer scraping: %v\n", err)
		s.stats.FailedScans++
	})
}

// ScrapeAlbumImages obtiene las primeras 3 imágenes de un álbum específico
func (s *YupooScraper) ScrapeAlbumImages(albumID string) ([]string, error) {
	images := make([]string, 0, imagesPerAlbum)

	// Crear un nuevo collector para este álbum específico
	c := colly.NewCollector(
		colly.AllowedDomains("huang-66.x.yupoo.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       500 * time.Millisecond, // Reducido de 2s a 0.5s
	})

	// Buscar las imágenes en el álbum
	c.OnHTML("div.image__main img", func(e *colly.HTMLElement) {
		if len(images) < imagesPerAlbum {
			imgURL := e.Attr("src")
			if imgURL != "" {
				// Convertir a URL completa si es relativa
				if strings.HasPrefix(imgURL, "//") {
					imgURL = "https:" + imgURL
				}
				// Cambiar "small" por "medium" para mejor calidad
				imgURL = strings.Replace(imgURL, "/small.jpg", "/medium.jpg", 1)
				images = append(images, imgURL)
			}
		}
	})

	url := fmt.Sprintf("%s/albums/%s?uid=1", baseURL, albumID)
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error al visitar álbum %s: %v", albumID, err)
	}

	c.Wait()

	return images, nil
}

// GetAlbums devuelve todos los álbumes scrapeados
func (s *YupooScraper) GetAlbums() []models.YupooAlbum {
	return s.albums
}

// GetStats devuelve las estadísticas del scraping
func (s *YupooScraper) GetStats() models.ScraperStats {
	return s.stats
}

// EnrichAlbumsWithImages obtiene las imágenes para cada álbum
func (s *YupooScraper) EnrichAlbumsWithImages() error {
	log.Println("\nObteniendo imágenes de cada álbum...")

	for i := range s.albums {
		log.Printf("  [%d/%d] Obteniendo imágenes del álbum: %s\n", i+1, len(s.albums), s.albums[i].Title)

		images, err := s.ScrapeAlbumImages(s.albums[i].ID)
		if err != nil {
			log.Printf("    ✗ Error: %v\n", err)
			continue
		}

		s.albums[i].Images = images
		log.Printf("    ✓ %d imágenes obtenidas\n", len(images))

		// Pequeña pausa entre álbumes
		time.Sleep(200 * time.Millisecond) // Reducido de 1s a 0.2s
	}

	return nil
}