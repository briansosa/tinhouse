package scraper

import (
	"context"
	"strings"

	"github.com/findhouse/internal/models"
	"github.com/findhouse/internal/scraper/tokko"
)

// PropertyScraper define la interfaz que deben implementar todos los scrapers de propiedades
type PropertyScraper interface {
	// SearchProperties busca propiedades en el sitio web de la inmobiliaria
	SearchProperties(ctx context.Context) ([]models.Property, error)

	// GetPropertyDetails obtiene los detalles de una propiedad específica
	GetPropertyDetails(ctx context.Context, url string) (*models.PropertyDetails, error)
}

// Verificación de que TokkoScraper implementa la interfaz PropertyScraper
var _ PropertyScraper = (*tokko.TokkoScraper)(nil)

// NewScraper crea un nuevo scraper basado en el sistema de la inmobiliaria
func NewScraper(sistema string, baseURL string) PropertyScraper {
	// Por ahora solo soportamos Tokko, pero en el futuro podemos agregar más
	sistema = strings.ToLower(sistema)

	if strings.Contains(sistema, "tokko") {
		return tokko.New(baseURL)
	}

	// Si no reconocemos el sistema, devolvemos nil
	// En el futuro podríamos devolver un scraper genérico o un error
	return nil
}
