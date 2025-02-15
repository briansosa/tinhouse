package db

import "time"

// Inmobiliaria representa una inmobiliaria en la base de datos
type Inmobiliaria struct {
	ID        int64     `db:"id"`
	Nombre    string    `db:"nombre"`
	URL       string    `db:"url"`
	Sistema   string    `db:"sistema"`
	Zona      string    `db:"zona"`
	Rating    float64   `db:"rating"`
	Direccion string    `db:"direccion"`
	Telefono  string    `db:"telefono"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Busqueda representa una búsqueda realizada
type Busqueda struct {
	ID           int64     `db:"id"`
	Operation    string    `db:"operation"`
	PropertyType string    `db:"property_type"`
	Zone         string    `db:"zone"`
	Location     string    `db:"location"`
	MinPriceUSD  float64   `db:"min_price_usd"`
	MaxPriceUSD  float64   `db:"max_price_usd"`
	CreatedAt    time.Time `db:"created_at"`
}

// Propiedad representa una propiedad en la base de datos
type Propiedad struct {
	ID             int64     `db:"id"`
	InmobiliariaID int64     `db:"inmobiliaria_id"`
	Codigo         string    `db:"codigo"`
	Titulo         string    `db:"titulo"`
	Precio         string    `db:"precio"`
	Direccion      string    `db:"direccion"`
	URL            string    `db:"url"`
	ImagenURL      string    `db:"imagen_url"`
	FechaScraping  time.Time `db:"fecha_scraping"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
}

// BusquedaPropiedad representa la relación entre búsquedas y propiedades
type BusquedaPropiedad struct {
	BusquedaID  int64     `db:"busqueda_id"`
	PropiedadID int64     `db:"propiedad_id"`
	CreatedAt   time.Time `db:"created_at"`
}
