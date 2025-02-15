package models

type Property struct {
	ID          string
	Title       string
	Type        string // casa, departamento, etc
	Operation   string // venta, alquiler
	PriceText   string // Cambiado de Price float64 a PriceText string
	Currency    string // USD, ARS
	Address     string
	Zone        string
	Bedrooms    int
	Bathrooms   int
	Area        string // Cambiado a string porque viene como "220 m²"
	Rooms       string // Nuevo campo para ambientes
	CoveredArea float64
	Agency      string
	URL         string
	Images      []string
	Code        string // Código de la propiedad
	ImageURL    string // URL de la imagen
}

type PropertyFilter struct {
	Operation   string
	Type        string
	Zone        string // e.g. "G.B.A. Zona Sur"
	Location    string // e.g. "Lanús"
	MinPriceUSD float64
	MaxPriceUSD float64
}
