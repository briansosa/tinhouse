package tokko

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/findhouse/internal/models"
)

type Scraper struct {
	BaseURL string
}

func New(baseURL string) *Scraper {
	cleanURL := strings.TrimRight(baseURL, "/")
	return &Scraper{
		BaseURL: cleanURL,
	}
}

func (s *Scraper) SearchProperties(ctx context.Context, filter models.PropertyFilter) ([]models.Property, error) {
	baseURL := strings.TrimRight(s.BaseURL, "/")

	url := fmt.Sprintf("%s/Buscar?operation=%s&ptypes=%s&locations=%s&min-price=%.0f&max-price=%.0f&currency=USD&o=2,2&l=1",
		baseURL,
		operationToTokko(filter.Operation),
		propertyTypeToTokko(filter.Type),
		locationToTokko(filter.Location),
		filter.MinPriceUSD,
		filter.MaxPriceUSD,
	)

	fmt.Printf("Buscando en URL: %s\n", url)

	var properties []models.Property

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`
			(() => {
				console.log('=== Debug Info ===');
				
				// Buscar el contenedor de propiedades
				const container = document.querySelector('#propiedades.resultados-list');
				console.log('Container found:', container ? 'yes' : 'no');
				
				// Buscar los items de propiedades (li elements dentro del contenedor)
				const properties = container ? Array.from(container.querySelectorAll('li')) : [];
				console.log('Properties found:', properties.length);
				
				return properties.map(prop => {
					// Extraer datos usando las clases que vemos en el DOM
					const priceElement = prop.querySelector('.prop-valor-nro');
					const codeElement = prop.querySelector('.codref');
					const titleElement = prop.querySelector('.prop-desc-tipo-ub');
					const addressElement = prop.querySelector('.prop-desc-dir');
					const imageElement = prop.querySelector('.dest-img');
					
					return {
						title: titleElement?.textContent?.trim(),
						priceText: priceElement?.textContent?.trim() || '',
						address: addressElement?.textContent?.trim(),
						code: codeElement?.textContent?.trim(),
						url: prop.querySelector('a')?.href,
						imageUrl: imageElement?.src || ''
					};
				});
			})()
		`, &properties),
	)

	if err != nil {
		return nil, fmt.Errorf("error extrayendo propiedades: %v", err)
	}

	return properties, nil
}

func operationToTokko(op string) string {
	switch op {
	case "venta":
		return "1"
	case "alquiler":
		return "2"
	default:
		return "1"
	}
}

func propertyTypeToTokko(t string) string {
	switch t {
	case "casa":
		return "3"
	case "departamento":
		return "2"
	case "terreno":
		return "1"
	case "local":
		return "6"
	case "oficina":
		return "7"
	case "galpon":
		return "8"
	default:
		return "3"
	}
}

// Nuevas funciones para el formato de URL
func propertyTypeToTokkoURL(t string) string {
	switch t {
	case "casa":
		return "Casa"
	case "departamento":
		return "Departamento"
	default:
		return "Casa"
	}
}

func operationToTokkoURL(op string) string {
	switch op {
	case "venta":
		return "Venta"
	case "alquiler":
		return "Alquiler"
	default:
		return "Venta"
	}
}

// Mapeo de ubicaciones
func locationToTokko(location string) string {
	locations := map[string]string{
		"Lanús": "26540",
		// Agregar más mappings según necesitemos
	}
	if id, ok := locations[location]; ok {
		return id
	}
	return "26540" // Default a Lanús
}
