# Ultima Linea - Yupoo Scraper

Scraper especializado para extraer información de productos de camisetas de fútbol desde Yupoo (huang-66.x.yupoo.com).

## Características

- Scraping automatizado por **categorías específicas de fútbol** (14 categorías)
- **Límite de 15 álbumes por categoría** para rapidez (configurable)
- **Descarga paralela optimizada** de imágenes (3 workers simultáneos)
- **Filtrado automático** de productos no deseados (NBA, NFL)
- Extracción de información de productos:
  - Título completo del producto
  - ID del álbum
  - Número de imágenes disponibles
  - URLs de las primeras 3 imágenes de cada producto
  - Categoría (Mundial, La Liga, Premier, Serie A, etc.)
  - Página de origen
- Rate limiting optimizado para velocidad
- Exportación de datos a JSON
- Estadísticas detalladas del proceso de scraping

## Categorías Incluidas

El scraper recorre únicamente las siguientes categorías de fútbol:

1. **Mundial 2026** - Camisetas del Mundial
2. **La Liga España** - Equipos españoles
3. **Liga Argentina** - Equipos argentinos
4. **Retro Collection** - Camisetas retro/clásicas
5. **Liga Chilena** - Equipos chilenos
6. **Brasileirao** - Liga brasileña
7. **Premier** - Premier League inglesa
8. **Serie A** - Liga italiana
9. **Ligue 1 Francia** - Liga francesa
10. **Bundesliga** - Liga alemana
11. **Saudi League** - Liga saudí
12. **MLS** - Major League Soccer
13. **World Jersey** - Selecciones nacionales
14. **Player Version** - Versiones de jugador

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

### Scraping básico (todas las categorías de fútbol)

```bash
go run cmd/scraper/main.go
```

Este comando scrapeará automáticamente las **14 categorías de fútbol** configuradas.

### Scraping con parámetros personalizados

```bash
# Scrapear sin imágenes (más rápido, solo metadata)
go run cmd/scraper/main.go -images=false

# Cambiar archivo de salida
go run cmd/scraper/main.go -output productos_futbol.json

# Scrapear con imágenes (default)
go run cmd/scraper/main.go -images=true
```

### Flags disponibles

| Flag | Tipo | Default | Descripción |
|------|------|---------|-------------|
| `-images` | bool | true | Obtener imágenes de cada álbum (3 por álbum) |
| `-output` | string | yupoo_products.json | Archivo de salida JSON |

**Nota**: Ya no hay flags `-start` y `-end` porque el scraper ahora trabaja con categorías específicas en lugar de páginas.

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

## Rendimiento y Optimizaciones

El scraper está optimizado para balance entre velocidad y respeto al servidor:

- **Delay entre requests**: 2 segundos (evita rate limiting)
- **Paralelismo de páginas**: 2 requests simultáneos
- **Descarga de imágenes**: 3 workers paralelos
- **Filtrado inteligente**: Excluye automáticamente NBA y NFL
- **User-Agent**: Mozilla/5.0 (simula navegador)

## Ejemplos de Uso

### Caso 1: Scraping rápido para testing (sin imágenes)

```bash
go run cmd/scraper/main.go -images=false
```

Esto scrapeará todas las categorías de fútbol sin obtener imágenes (15 álbumes por categoría), **ultra rápido** para pruebas.

### Caso 2: Scraping completo con imágenes

```bash
go run cmd/scraper/main.go -output catalogo_futbol.json
```

Esto scrapeará las 14 categorías con imágenes (15 álbumes por categoría = 210 productos totales). Gracias a la descarga paralela (3 workers), es más rápido sin saturar el servidor.

### Caso 3: Integración con CI/CD

```bash
# Ejecución automática diaria
go run cmd/scraper/main.go -output "futbol_$(date +%Y%m%d).json"
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

- ✅ **Optimizado para estabilidad**: Descarga paralela con 3 workers y rate limiting respetuoso
- ✅ **Solo fútbol**: Filtra automáticamente NBA, NFL y otros deportes
- ✅ **14 categorías específicas**: Solo las ligas y torneos relevantes (15 álbumes por categoría)
- ⚙️ **Configurable**: Puedes cambiar `albumsPerCategory` en `internal/scraper/yupoo.go` para scrapear más o menos álbumes
- Las imágenes se obtienen en calidad "medium" (mejor que "small", más rápido que "original")
- El scraping incluye delays de 2 segundos entre requests para evitar bloqueos del servidor
- Algunos álbumes pueden estar protegidos con contraseña y no serán accesibles

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