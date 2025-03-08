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
	TipoPropiedad      *int64   `db:"tipo_propiedad"`
	Ubicacion          *string  `db:"ubicacion"`
	Dormitorios        *int     `db:"dormitorios"`
	Banios             *int     `db:"banios"`
	Antiguedad         *int     `db:"antiguedad"`
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
	Latitud            *float64 `db:"latitud"`
	Longitud           *float64 `db:"longitud"`

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
	PropertyTypeID    *int64   `json:"property_type_id"`
	PropertyTypeIDs   []int64  `json:"property_type_ids"`
	Locations         []string `json:"locations"`
	Features          []string `json:"features"`
	PriceMin          *float64 `json:"price_min"`
	PriceMax          *float64 `json:"price_max"`
	Currency          string   `json:"currency"`
	SizeMin           *float64 `json:"size_min"`         // Compatibilidad con versión anterior
	SizeMax           *float64 `json:"size_max"`         // Compatibilidad con versión anterior
	TotalAreaMin      *float64 `json:"total_area_min"`   // Superficie Total mínima
	TotalAreaMax      *float64 `json:"total_area_max"`   // Superficie Total máxima
	CoveredAreaMin    *float64 `json:"covered_area_min"` // Superficie Cubierta mínima
	CoveredAreaMax    *float64 `json:"covered_area_max"` // Superficie Cubierta máxima
	LandAreaMin       *float64 `json:"land_area_min"`    // Superficie Terreno mínima
	LandAreaMax       *float64 `json:"land_area_max"`    // Superficie Terreno máxima
	Front             *float64 `json:"front"`            // Frente
	Back              *float64 `json:"back"`             // Fondo
	Rooms             *int     `json:"rooms"`
	Bathrooms         *int     `json:"bathrooms"`
	Antiquity         *int     `json:"antiquity"`
	Disposition       []string `json:"disposition"`    // Disposición de la propiedad
	Orientation       []string `json:"orientation"`    // Orientación de la propiedad
	Condition         []string `json:"condition"`      // Condición de la propiedad
	OperationType     []string `json:"operation_type"` // Tipo de operación de la propiedad
	Situation         []string `json:"situation"`      // Situación de la propiedad
	AgencyIDs         []int64  `json:"agencies"`       // IDs de inmobiliarias
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

// PropertyType representa un tipo de propiedad normalizado
type PropertyType struct {
	ID        int64     `db:"id" json:"id"`
	Code      string    `db:"code" json:"code"`
	Name      string    `db:"name" json:"name"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// Lista genérica para valores predefinidos
type List struct {
	ID          int64     `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// Valor de una lista
type ListValue struct {
	ID          int64     `db:"id" json:"id"`
	ListID      int64     `db:"list_id" json:"list_id"`
	Value       string    `db:"value" json:"value"`
	DisplayName string    `db:"display_name" json:"display_name"`
	SortOrder   int       `db:"sort_order" json:"sort_order"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
