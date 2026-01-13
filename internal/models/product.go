package models

import "time"

// YupooAlbum representa un álbum de Yupoo con información del producto
type YupooAlbum struct {
	ID          string    `json:"id"`           // ID del álbum en Yupoo
	Title       string    `json:"title"`        // Título completo del producto
	ImageCount  int       `json:"image_count"`  // Número de imágenes en el álbum
	Images      []string  `json:"images"`       // URLs de las primeras 3 imágenes
	Category    string    `json:"category"`     // Categoría del producto
	PageNumber  int       `json:"page_number"`  // Página donde fue encontrado
	AlbumURL    string    `json:"album_url"`    // URL completa del álbum
	ScrapedAt   time.Time `json:"scraped_at"`   // Fecha y hora del scraping
}

// ScraperStats representa estadísticas del proceso de scraping
type ScraperStats struct {
	TotalPages      int       `json:"total_pages"`
	TotalAlbums     int       `json:"total_albums"`
	SuccessfulScans int       `json:"successful_scans"`
	FailedScans     int       `json:"failed_scans"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Duration        string    `json:"duration"`
}