package scraper

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/ultima-linea/scraper/internal/models"
)

const (
	baseURL           = "https://huang-66.x.yupoo.com"
	imagesPerAlbum    = 3  // Solo las primeras 3 imágenes
	albumsPerCategory = 15 // Solo los primeros 15 álbumes por categoría
)

// Category representa una categoría de productos
type Category struct {
	ID   string
	Name string
	URL  string
}

// FootballCategories son las categorías de fútbol a scrapear
var FootballCategories = []Category{
	{ID: "661649", Name: "Mundial 2026", URL: baseURL + "/categories/661649"},
	{ID: "661476", Name: "La Liga España", URL: baseURL + "/categories/661476"},
	{ID: "3925870", Name: "Liga Argentina", URL: baseURL + "/categories/3925870"},
	{ID: "3258137", Name: "Retro Collection", URL: baseURL + "/categories/3258137"},
	{ID: "3185811", Name: "Liga Chilena", URL: baseURL + "/categories/3185811"},
	{ID: "3534092", Name: "Brasileirao", URL: baseURL + "/categories/3534092"},
	{ID: "654538", Name: "Premier", URL: baseURL + "/categories/654538"},
	{ID: "654560", Name: "Serie A", URL: baseURL + "/categories/654560"},
	{ID: "661562", Name: "Ligue 1 Francia", URL: baseURL + "/categories/661562"},
	{ID: "660913", Name: "Bundesliga", URL: baseURL + "/categories/660913"},
	{ID: "4279143", Name: "Saudi League", URL: baseURL + "/categories/4279143"},
	{ID: "3925867", Name: "MLS", URL: baseURL + "/categories/3925867"},
	{ID: "3576323", Name: "World Jersey", URL: baseURL + "/categories/3576323"},
	{ID: "3341199", Name: "Player Version", URL: baseURL + "/categories/3341199"},
}

type YupooScraper struct {
	collector           *colly.Collector
	albums              []models.YupooAlbum
	stats               models.ScraperStats
	currentCategory     string // Para trackear la categoría actual durante el scraping
	albumsInCategory    int    // Contador de álbumes en la categoría actual
	mu                  sync.Mutex // Para proteger acceso concurrente a albums
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
		Parallelism: 2,        // Reducido a 2 para evitar rate limiting
		Delay:       2 * time.Second, // Aumentado a 2s para ser más respetuoso
	})

	return &YupooScraper{
		collector: c,
		albums:    make([]models.YupooAlbum, 0),
		stats: models.ScraperStats{
			StartTime: time.Now(),
		},
	}
}

// ScrapeCategories recorre todas las categorías de fútbol especificadas
func (s *YupooScraper) ScrapeCategories(categories []Category) error {
	log.Printf("Iniciando scraping de %d categorías de fútbol...\n", len(categories))

	// Configurar el scraper para el listado de álbumes
	s.setupAlbumListingCollector()

	// Recorrer cada categoría
	for i, category := range categories {
		log.Printf("\n[%d/%d] Scraping categoría: %s\n", i+1, len(categories), category.Name)

		// Establecer la categoría actual y resetear contador
		s.currentCategory = category.Name
		s.albumsInCategory = 0

		// Cada categoría puede tener múltiples páginas
		// Empezamos por la página 1 y vemos si hay más
		pageNum := 1
		hasMorePages := true

		for hasMorePages {
			var url string
			if pageNum == 1 {
				url = category.URL
			} else {
				url = fmt.Sprintf("%s?page=%d", category.URL, pageNum)
			}

			log.Printf("  Scraping página %d: %s\n", pageNum, url)

			// Crear un collector temporal para esta categoría
			tempAlbumsCount := len(s.albums)

			err := s.collector.Visit(url)
			if err != nil {
				log.Printf("  Error al visitar página %d: %v\n", pageNum, err)
				s.stats.FailedScans++
				break
			}

			s.collector.Wait()

			// Si no encontramos nuevos álbumes, asumimos que no hay más páginas
			if len(s.albums) == tempAlbumsCount {
				hasMorePages = false
			} else {
				pageNum++
				s.stats.TotalPages++
				// Limitar a 15 álbumes por categoría
				if s.albumsInCategory >= albumsPerCategory {
					log.Printf("  Alcanzado límite de %d álbumes para esta categoría\n", albumsPerCategory)
					hasMorePages = false
				}
				// Limitar a un número razonable de páginas por categoría
				if pageNum > 50 {
					log.Printf("  Alcanzado límite de 50 páginas para esta categoría\n")
					hasMorePages = false
				}
			}
		}

		log.Printf("  ✓ Categoría %s completada: %d álbumes en esta categoría\n", category.Name, s.albumsInCategory)
	}

	s.stats.EndTime = time.Now()
	s.stats.Duration = s.stats.EndTime.Sub(s.stats.StartTime).String()
	s.stats.TotalAlbums = len(s.albums)

	log.Printf("\nScraping de categorías completado: %d álbumes encontrados en total\n", s.stats.TotalAlbums)

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

		// Filtrar álbumes no deseados (NBA, NFL)
		titleUpper := strings.ToUpper(album.Title)
		if strings.Contains(titleUpper, "NBA") || strings.Contains(titleUpper, "NFL") {
			return // Ignorar este álbum
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

		// Asignar la categoría actual
		album.Category = s.currentCategory

		// Solo agregar si tiene datos válidos
		if album.ID != "" && album.Title != "" {
			s.mu.Lock()
			// Verificar si ya alcanzamos el límite para esta categoría
			if s.albumsInCategory < albumsPerCategory {
				s.albums = append(s.albums, album)
				s.stats.SuccessfulScans++
				s.albumsInCategory++
				s.mu.Unlock()
				log.Printf("  ✓ Álbum encontrado: %s (ID: %s, %d imágenes)\n", album.Title, album.ID, album.ImageCount)
			} else {
				s.mu.Unlock()
			}
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
		Delay:       2 * time.Second, // Aumentado a 2s para evitar rate limiting
	})

	// Buscar las imágenes en el álbum
	// Las imágenes en Yupoo están con este patrón: img[src*="photo.yupoo.com"]
	c.OnHTML("img[src*='photo.yupoo.com/huang-66']", func(e *colly.HTMLElement) {
		if len(images) < imagesPerAlbum {
			imgURL := e.Attr("src")
			if imgURL != "" && !strings.Contains(imgURL, "logo") && !strings.Contains(imgURL, "icon") {
				// Convertir a URL completa si es relativa (formato //photo.yupoo.com/...)
				if strings.HasPrefix(imgURL, "//") {
					imgURL = "https:" + imgURL
				}
				// Convertir de small a medium para mejor calidad
				imgURL = strings.Replace(imgURL, "/small.jpg", "/medium.jpg", 1)
				imgURL = strings.Replace(imgURL, "/small.jpeg", "/medium.jpeg", 1)
				images = append(images, imgURL)
			}
		}
	})

	// Usar el formato correcto de URL: /albums/{albumID}?uid=1
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

// EnrichAlbumsWithImages obtiene las imágenes para cada álbum usando concurrencia
func (s *YupooScraper) EnrichAlbumsWithImages() error {
	log.Println("\nObteniendo imágenes de cada álbum...")

	// Usar workers paralelos para acelerar la descarga
	const maxWorkers = 3 // Reducido de 20 a 3 para evitar rate limiting
	semaphore := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex // Para proteger el contador de logs

	for i := range s.albums {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			semaphore <- struct{}{}        // Adquirir slot
			defer func() { <-semaphore }() // Liberar slot

			mu.Lock()
			log.Printf("  [%d/%d] Obteniendo imágenes del álbum: %s\n", index+1, len(s.albums), s.albums[index].Title)
			mu.Unlock()

			images, err := s.ScrapeAlbumImages(s.albums[index].ID)
			if err != nil {
				mu.Lock()
				log.Printf("    ✗ Error: %v\n", err)
				mu.Unlock()
				return
			}

			s.albums[index].Images = images
			mu.Lock()
			log.Printf("    ✓ %d imágenes obtenidas\n", len(images))
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	return nil
}