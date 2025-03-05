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
	Moneda         string    `db:"moneda"`
	Direccion      string    `db:"direccion"`
	URL            string    `db:"url"`
	ImagenURL      string    `db:"imagen_url"`
	Imagenes       *[]string `db:"imagenes"`
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
	Operacion          *string  `db:"operacion"`
	Condicion          *string  `db:"condicion"`
	Orientacion        *string  `db:"orientacion"`
	Disposicion        *string  `db:"disposicion"`

	// Campo virtual para indicar si es favorita
	IsFavorite bool `db:"-"`

	// Campo virtual para características agrupadas por categoría
	Features map[string][]string `db:"-"`
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
	IsFavorite bool      `db:"is_favorite"`
	CreatedAt  time.Time `db:"created_at"`
}

// PropertyNote representa una nota de propiedad en la base de datos
type PropertyNote struct {
	ID         int64     `db:"id" json:"id"`
	PropertyID int64     `db:"property_id" json:"property_id"`
	Text       string    `db:"note" json:"text"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
}

// PropertyFilter representa los filtros aplicables a las propiedades
type PropertyFilter struct {
	PropertyType      string   `json:"property_type"`
	Locations         []string `json:"locations"`
	Features          []string `json:"features"`
	PriceMin          *float64 `json:"price_min"`
	PriceMax          *float64 `json:"price_max"`
	Currency          string   `json:"currency"`
	SizeMin           *float64 `json:"size_min"`
	SizeMax           *float64 `json:"size_max"`
	Rooms             *int     `json:"rooms"`
	Bathrooms         *int     `json:"bathrooms"`
	Antiquity         *int     `json:"antiquity"`
	ShowOnlyWithNotes bool     `json:"show_only_with_notes"`
	ShowOnlyFavorites bool     `json:"show_only_favorites"`
}

// PropertyFeature representa una característica de una propiedad
type PropertyFeature struct {
	ID        int64     `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Category  string    `db:"category" json:"category"` // 'servicio', 'ambiente', 'adicional'
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// PropertyFeatureRelation representa la relación entre una propiedad y una característica
type PropertyFeatureRelation struct {
	PropertyID int64     `db:"property_id" json:"property_id"`
	FeatureID  int64     `db:"feature_id" json:"feature_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
