package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ultima-linea/scraper/internal/scraper"
	"github.com/ultima-linea/scraper/pkg/utils"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando valores por defecto")
	}

	// Flags de línea de comandos
	withImages := flag.Bool("images", true, "Obtener imágenes de cada álbum")
	outputFile := flag.String("output", "yupoo_products.json", "Archivo de salida JSON")
	flag.Parse()

	// Banner
	fmt.Println("╔════════════════════════════════════════════╗")
	fmt.Println("║   Ultima Linea - Yupoo Scraper            ║")
	fmt.Println("║   Scraping de categorías de fútbol        ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Println()

	// Crear scraper
	yupooScraper := scraper.NewYupooScraper()

	// Scraping de categorías de fútbol
	log.Printf("Configuración: %d categorías, Imágenes: %v\n", len(scraper.FootballCategories), *withImages)
	fmt.Println()

	err := yupooScraper.ScrapeCategories(scraper.FootballCategories)
	if err != nil {
		log.Fatalf("Error durante el scraping: %v\n", err)
	}

	// Obtener imágenes si está habilitado
	if *withImages {
		err = yupooScraper.EnrichAlbumsWithImages()
		if err != nil {
			log.Printf("Advertencia: Error al obtener imágenes: %v\n", err)
		}
	}

	// Obtener resultados
	albums := yupooScraper.GetAlbums()
	stats := yupooScraper.GetStats()

	// Guardar a JSON
	log.Printf("\nGuardando resultados en %s...\n", *outputFile)
	err = utils.SaveToJSON(*outputFile, map[string]interface{}{
		"stats":  stats,
		"albums": albums,
	})
	if err != nil {
		log.Fatalf("Error al guardar JSON: %v\n", err)
	}

	// Mostrar estadísticas finales
	fmt.Println("\n╔════════════════════════════════════════════╗")
	fmt.Println("║           ESTADÍSTICAS FINALES             ║")
	fmt.Println("╚════════════════════════════════════════════╝")
	fmt.Printf("📄 Total de páginas scrapeadas:  %d\n", stats.TotalPages)
	fmt.Printf("📦 Total de álbumes encontrados: %d\n", stats.TotalAlbums)
	fmt.Printf("✅ Scraping exitosos:            %d\n", stats.SuccessfulScans)
	fmt.Printf("❌ Scraping fallidos:            %d\n", stats.FailedScans)
	fmt.Printf("⏱️  Duración total:               %s\n", stats.Duration)
	fmt.Printf("💾 Archivo de salida:            %s\n", *outputFile)
	fmt.Println()

	// Verificar que el archivo existe
	if _, err := os.Stat(*outputFile); err == nil {
		fileInfo, _ := os.Stat(*outputFile)
		fmt.Printf("✓ Archivo guardado exitosamente (%.2f KB)\n", float64(fileInfo.Size())/1024)
	}

	fmt.Println("\n✨ Scraping completado exitosamente!")
}