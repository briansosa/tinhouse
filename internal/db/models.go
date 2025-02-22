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

	// Campos que pueden ser NULL
	TipoPropiedad      *string  `db:"tipo_propiedad"`
	Ubicacion          *string  `db:"ubicacion"`
	Dormitorios        *int     `db:"dormitorios"`
	Banios             *int     `db:"banios"`
	Antiguedad         *string  `db:"antiguedad"`
	SuperficieCubierta *float64 `db:"superficie_cubierta"`
	SuperficieTotal    *float64 `db:"superficie_total"`
	SuperficieTerreno  *float64 `db:"superficie_terreno"`
	Frente             *float64 `db:"frente"`
	Fondo              *float64 `db:"fondo"`
	Ambientes          *int     `db:"ambientes"`
	Plantas            *int     `db:"plantas"`
	Cocheras           *int     `db:"cocheras"`
	Situacion          *string  `db:"situacion"`
	Expensas           *float64 `db:"expensas"`
	Descripcion        *string  `db:"descripcion"`
	Status             string   `db:"status"`
}

// BusquedaPropiedad representa la relación entre búsquedas y propiedades
type BusquedaPropiedad struct {
	BusquedaID  int64     `db:"busqueda_id"`
	PropiedadID int64     `db:"propiedad_id"`
	CreatedAt   time.Time `db:"created_at"`
}

// PropertyRating representa una calificación de propiedad en la base de datos
type PropertyRating struct {
	ID         int64     `db:"id"`
	PropertyID int64     `db:"property_id"`
	Rating     string    `db:"rating"` // 'like' o 'dislike'
	CreatedAt  time.Time `db:"created_at"`
}
