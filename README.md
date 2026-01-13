# Ultima Linea - Yupoo Scraper

Scraper especializado para extraer información de productos de camisetas desde Yupoo (huang-66.x.yupoo.com).

## Características

- Scraping automatizado de todas las páginas de la galería (45 páginas)
- Extracción de información de productos:
  - Título completo del producto
  - ID del álbum
  - Número de imágenes disponibles
  - URLs de las primeras 3 imágenes de cada producto
  - Categoría y página de origen
- Rate limiting configurable para no saturar el servidor
- Exportación de datos a JSON
- Estadísticas detalladas del proceso de scraping

## Tecnologías

- Go 1.24+
- Colly v2 (Web Scraping Framework)
- godotenv (Variables de entorno)

## Estructura del Proyecto

```
.
├── cmd/
│   └── scraper/
│       └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── scraper/
│   │   └── yupoo.go            # Lógica de scraping de Yupoo
│   └── models/
│       └── product.go          # Modelos de datos
├── pkg/
│   └── utils/
│       └── json.go             # Utilidades para JSON
├── .env                        # Configuración del scraper
├── .gitignore
├── go.mod
└── README.md
```

## Instalación

### 1. Clonar el repositorio

```bash
git clone <repository-url>
cd ultima-linea-scraper
```

### 2. Instalar dependencias

```bash
go mod download
```

### 3. Configurar variables de entorno (Opcional)

El archivo `.env` ya viene con valores por defecto:

```env
SCRAPER_START_PAGE=1
SCRAPER_END_PAGE=45
SCRAPER_WITH_IMAGES=true
SCRAPER_OUTPUT_FILE=yupoo_products.json
SCRAPER_DELAY=2
SCRAPER_PARALLELISM=2
```

## Uso

### Scraping básico (primeras 5 páginas)

```bash
go run cmd/scraper/main.go
```

### Scraping con parámetros personalizados

```bash
# Scrapear páginas 1 a 10
go run cmd/scraper/main.go -start 1 -end 10

# Scrapear sin imágenes (más rápido)
go run cmd/scraper/main.go -start 1 -end 10 -images=false

# Cambiar archivo de salida
go run cmd/scraper/main.go -output productos.json

# Scrapear todas las páginas
go run cmd/scraper/main.go -start 1 -end 45
```

### Flags disponibles

| Flag | Tipo | Default | Descripción |
|------|------|---------|-------------|
| `-start` | int | 1 | Página inicial para scraping |
| `-end` | int | 5 | Página final para scraping |
| `-images` | bool | true | Obtener imágenes de cada álbum |
| `-output` | string | yupoo_products.json | Archivo de salida JSON |

### Compilar binario

```bash
# Compilar para Windows
go build -o scraper.exe cmd/scraper/main.go

# Compilar para Linux
GOOS=linux GOARCH=amd64 go build -o scraper cmd/scraper/main.go

# Ejecutar el binario
./scraper.exe -start 1 -end 45
```

## Formato de Salida

El scraper genera un archivo JSON con la siguiente estructura:

```json
{
  "stats": {
    "total_pages": 5,
    "total_albums": 120,
    "successful_scans": 120,
    "failed_scans": 0,
    "start_time": "2026-01-13T13:00:00Z",
    "end_time": "2026-01-13T13:15:00Z",
    "duration": "15m0s"
  },
  "albums": [
    {
      "id": "223106510",
      "title": "2026 Brazil Hollywood Keeper Long sleeved Fan version",
      "image_count": 21,
      "images": [
        "https://photo.yupoo.com/huang-66/abc123/medium.jpg",
        "https://photo.yupoo.com/huang-66/def456/medium.jpg",
        "https://photo.yupoo.com/huang-66/ghi789/medium.jpg"
      ],
      "category": "",
      "page_number": 1,
      "album_url": "https://huang-66.x.yupoo.com/albums/223106510",
      "scraped_at": "2026-01-13T13:05:00Z"
    }
  ]
}
```

## Rate Limiting

El scraper está configurado con rate limiting para evitar saturar el servidor:

- **Delay entre requests**: 2 segundos (configurable)
- **Paralelismo**: 2 requests simultáneos (configurable)
- **User-Agent**: Mozilla/5.0 (simula navegador)

## Ejemplos de Uso

### Caso 1: Scraping rápido para testing (sin imágenes)

```bash
go run cmd/scraper/main.go -start 1 -end 2 -images=false
```

Esto scrapeará solo las primeras 2 páginas sin obtener imágenes, ideal para pruebas rápidas.

### Caso 2: Scraping completo de todo el catálogo

```bash
go run cmd/scraper/main.go -start 1 -end 45 -output catalogo_completo.json
```

Esto scrapeará las 45 páginas completas con imágenes. Puede tomar 30-60 minutos dependiendo de la conexión.

### Caso 3: Scraping incremental

```bash
# Primera ejecución: páginas 1-15
go run cmd/scraper/main.go -start 1 -end 15 -output parte1.json

# Segunda ejecución: páginas 16-30
go run cmd/scraper/main.go -start 16 -end 30 -output parte2.json

# Tercera ejecución: páginas 31-45
go run cmd/scraper/main.go -start 31 -end 45 -output parte3.json
```

## Desarrollo

### Formato de código

```bash
go fmt ./...
```

### Tests (próximamente)

```bash
go test ./...
```

## Notas Importantes

- El scraper respeta rate limits para no saturar el servidor de Yupoo
- Algunos álbumes pueden estar protegidos con contraseña y no serán accesibles
- Las imágenes se obtienen en calidad "medium" (mejor que "small", más rápido que "original")
- El scraping completo de 45 páginas puede tomar tiempo considerable
- Se recomienda hacer scraping incremental para catálogos grandes

## Información del Proveedor

- **Tienda Yupoo**: huang-66.x.yupoo.com
- **WhatsApp**: +86 13535386151
- **Total de páginas**: 45
- **Productos aproximados**: 1000+ álbumes

## Próximas Mejoras

- [ ] Integración directa con la base de datos del backend
- [ ] Categorización automática de productos
- [ ] Extracción de tallas desde el título
- [ ] Parseo inteligente de equipos y ligas
- [ ] Sistema de cron jobs para actualización automática
- [ ] Detección de productos nuevos
- [ ] Soporte para múltiples proveedores de Yupoo

## Licencia

[Especifica tu licencia aquí]