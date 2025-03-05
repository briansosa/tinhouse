package db

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/findhouse/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Leer y ejecutar el schema
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return nil, err
	}

	database := &DB{db}

	return database, nil
}

// CreateInmobiliaria inserta una nueva inmobiliaria en la base de datos
func (db *DB) CreateInmobiliaria(i *Inmobiliaria) error {
	query := `
		INSERT INTO inmobiliarias (nombre, url, sistema, zona, rating, direccion, telefono)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at`

	return db.QueryRow(query,
		i.Nombre, i.URL, i.Sistema, i.Zona, i.Rating, i.Direccion, i.Telefono,
	).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
}

// CreateBusqueda inserta una nueva búsqueda en la base de datos
func (db *DB) CreateBusqueda(b *Busqueda) error {
	query := `
        INSERT INTO busquedas (operation, property_type, zone, location, min_price_usd, max_price_usd)
        VALUES (?, ?, ?, ?, ?, ?)
        RETURNING id, created_at`

	return db.QueryRow(query,
		b.Operation, b.PropertyType, b.Zone, b.Location, b.MinPriceUSD, b.MaxPriceUSD,
	).Scan(&b.ID, &b.CreatedAt)
}

// CreatePropiedad inserta una nueva propiedad en la base de datos
func (db *DB) CreatePropiedad(p *Propiedad) error {
	query := `
        INSERT INTO propiedades (
            inmobiliaria_id, codigo, titulo, precio, moneda, direccion, url, imagen_url,
            tipo_propiedad, ubicacion, dormitorios, banios, antiguedad, 
            superficie_cubierta, superficie_total, frente, fondo, ambientes,
            expensas, descripcion, status
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(codigo) DO UPDATE SET
            titulo = excluded.titulo,
            precio = excluded.precio,
            moneda = excluded.moneda,
            direccion = excluded.direccion,
            url = excluded.url,
            imagen_url = excluded.imagen_url,
            tipo_propiedad = excluded.tipo_propiedad,
            ubicacion = excluded.ubicacion,
            dormitorios = excluded.dormitorios,
            banios = excluded.banios,
            antiguedad = excluded.antiguedad,
            superficie_cubierta = excluded.superficie_cubierta,
            superficie_total = excluded.superficie_total,
            frente = excluded.frente,
            fondo = excluded.fondo,
            ambientes = excluded.ambientes,
            expensas = excluded.expensas,
            descripcion = excluded.descripcion,
            status = excluded.status,
            updated_at = CURRENT_TIMESTAMP
        RETURNING id, created_at, updated_at`

	return db.QueryRow(query,
		p.InmobiliariaID, p.Codigo, p.Titulo, p.Precio, p.Moneda, p.Direccion, p.URL, p.ImagenURL,
		p.TipoPropiedad, p.Ubicacion, p.Dormitorios, p.Banios, p.Antiguedad,
		p.SuperficieCubierta, p.SuperficieTotal, p.Frente, p.Fondo, p.Ambientes,
		p.Expensas, p.Descripcion, p.Status,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

// LinkBusquedaPropiedad vincula una búsqueda con una propiedad
func (db *DB) LinkBusquedaPropiedad(busquedaID, propiedadID int64) error {
	query := `
        INSERT INTO busquedas_propiedades (busqueda_id, propiedad_id)
        VALUES (?, ?)
        ON CONFLICT DO NOTHING`

	_, err := db.Exec(query, busquedaID, propiedadID)
	return err
}

// GetInmobiliariasSinSistema retorna las inmobiliarias que no tienen sistema identificado
func (db *DB) GetInmobiliariasSinSistema() ([]Inmobiliaria, error) {
	query := `
		SELECT id, nombre, url, sistema, zona, rating, direccion, telefono, created_at, updated_at
		FROM inmobiliarias
		WHERE sistema IS NULL OR sistema = '' or sistema = 'No identificado'`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inmobiliarias []Inmobiliaria
	for rows.Next() {
		var i Inmobiliaria
		err := rows.Scan(
			&i.ID, &i.Nombre, &i.URL, &i.Sistema, &i.Zona, &i.Rating,
			&i.Direccion, &i.Telefono, &i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inmobiliarias = append(inmobiliarias, i)
	}

	return inmobiliarias, nil
}

// UpdateInmobiliariaSistema actualiza el sistema de una inmobiliaria
func (db *DB) UpdateInmobiliariaSistema(i *Inmobiliaria) error {
	query := `
		UPDATE inmobiliarias
		SET sistema = ?, updated_at = ?
		WHERE id = ?`

	_, err := db.Exec(query, i.Sistema, i.UpdatedAt, i.ID)
	return err
}

// normalizarTexto normaliza el texto para comparación
func normalizarTexto(texto string) string {
	texto = strings.ToLower(texto)

	// Reemplazar los mismos caracteres que en SQL
	replacer := strings.NewReplacer(
		" ", "",
		".", "",
		",", "",
		"-", "",
		"&", "",
		"'", "",
	)
	return replacer.Replace(texto)
}

// ExistsInmobiliaria verifica si ya existe una inmobiliaria similar
func (db *DB) ExistsInmobiliaria(nombre, direccion string) (bool, error) {
	var exists bool

	nombreNormalizado := normalizarTexto(nombre)
	direccionNormalizada := normalizarTexto(direccion)

	query := `
		WITH normalized_data AS (
			SELECT 
				id,
				LOWER(
					REPLACE(
						REPLACE(
							REPLACE(
								REPLACE(
									REPLACE(
										REPLACE(nombre, ' ', ''),
										'.', ''
									),
									',', ''
								),
								'-', ''
							),
							'&', ''
						),
						'''', ''
					)
				) as nombre_norm,
				LOWER(
					REPLACE(
						REPLACE(
							REPLACE(
								REPLACE(
									REPLACE(
										REPLACE(direccion, ' ', ''),
										'.', ''
									),
									',', ''
								),
								'-', ''
							),
							'&', ''
						),
						'''', ''
					)
				) as direccion_norm
			FROM inmobiliarias
		)
		SELECT EXISTS(
			SELECT 1 FROM normalized_data 
			WHERE 
				nombre_norm LIKE ? 
				OR (
					direccion_norm != ''
					AND direccion_norm LIKE ?
				)
		)`

	// Agregar comodines para búsqueda más flexible
	nombreBusqueda := "%" + nombreNormalizado + "%"
	direccionBusqueda := "%" + direccionNormalizada + "%"
	err := db.QueryRow(query, nombreBusqueda, direccionBusqueda).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetInmobiliariasSistema retorna las inmobiliarias que tienen sistema identificado
func (db *DB) GetInmobiliariasSistema() ([]Inmobiliaria, error) {
	query := `
		SELECT id, nombre, url, sistema, zona, rating, direccion, telefono, created_at, updated_at
		FROM inmobiliarias
		WHERE sistema IS NOT NULL 
		AND sistema != '' 
		AND sistema != 'No identificado'
		ORDER BY nombre`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inmobiliarias []Inmobiliaria
	for rows.Next() {
		var i Inmobiliaria
		err := rows.Scan(
			&i.ID, &i.Nombre, &i.URL, &i.Sistema, &i.Zona, &i.Rating,
			&i.Direccion, &i.Telefono, &i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		inmobiliarias = append(inmobiliarias, i)
	}

	return inmobiliarias, nil
}

// GetOrCreateBusqueda verifica si existe una búsqueda con los mismos parámetros
// Si existe, la retorna. Si no existe, la crea.
func (db *DB) GetOrCreateBusqueda(filter models.PropertyFilter) (*Busqueda, error) {
	query := `
		SELECT id, created_at 
		FROM busquedas 
		WHERE operation = ? 
		AND property_type = ? 
		AND zone = ? 
		AND location = ? 
		AND min_price_usd = ? 
		AND max_price_usd = ?`

	var busqueda Busqueda
	err := db.QueryRow(query,
		filter.Operation,
		filter.Type,
		filter.Zone,
		filter.Location,
		filter.MinPriceUSD,
		filter.MaxPriceUSD,
	).Scan(&busqueda.ID, &busqueda.CreatedAt)

	if err == sql.ErrNoRows {
		// La búsqueda no existe, la creamos
		busqueda = Busqueda{
			Operation:    filter.Operation,
			PropertyType: filter.Type,
			Zone:         filter.Zone,
			Location:     filter.Location,
			MinPriceUSD:  filter.MinPriceUSD,
			MaxPriceUSD:  filter.MaxPriceUSD,
		}

		if err := db.CreateBusqueda(&busqueda); err != nil {
			return nil, fmt.Errorf("error creando búsqueda: %v", err)
		}
		return &busqueda, nil
	}

	if err != nil {
		return nil, fmt.Errorf("error verificando búsqueda existente: %v", err)
	}

	return &busqueda, nil
}

// CreatePropiedadAndLink crea una propiedad si no existe y la vincula con una búsqueda
func (db *DB) CreatePropiedadAndLink(p *Propiedad, busquedaID int64) error {
	// Primero intentamos crear/actualizar la propiedad
	err := db.CreatePropiedad(p)
	if err != nil {
		return fmt.Errorf("error creando/actualizando propiedad: %v", err)
	}

	// Luego creamos el vínculo con la búsqueda
	err = db.LinkBusquedaPropiedad(busquedaID, p.ID)
	if err != nil {
		return fmt.Errorf("error vinculando búsqueda-propiedad: %v", err)
	}

	return nil
}

// GetPropiedadesSinDetalles retorna las propiedades que no tienen detalles completos
func (db *DB) GetPropiedadesSinDetalles() ([]Propiedad, error) {
	query := `
		SELECT 
			id, inmobiliaria_id, codigo, titulo, precio, direccion, url, imagen_url,
			NULLIF(imagenes, '') as imagenes,
			created_at, updated_at,
			tipo_propiedad, ubicacion, 
			NULLIF(dormitorios, '') as dormitorios,
			NULLIF(banios, '') as banios,
			antiguedad,
			NULLIF(superficie_cubierta, '') as superficie_cubierta,
			NULLIF(superficie_total, '') as superficie_total,
			NULLIF(superficie_terreno, '') as superficie_terreno,
			NULLIF(frente, '') as frente,
			NULLIF(fondo, '') as fondo,
			NULLIF(ambientes, '') as ambientes,
			NULLIF(plantas, '') as plantas,
			NULLIF(cocheras, '') as cocheras,
			situacion,
			NULLIF(expensas, '') as expensas,
			descripcion, status, operacion, condicion, orientacion, disposicion
		FROM propiedades
		WHERE status = 'pending'
		ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando propiedades sin detalles: %v", err)
	}
	defer rows.Close()

	var propiedades []Propiedad
	for rows.Next() {
		var p Propiedad
		var imagenesJSON sql.NullString // Para manejar NULL en la base de datos

		err := rows.Scan(
			&p.ID, &p.InmobiliariaID, &p.Codigo, &p.Titulo, &p.Precio, &p.Direccion,
			&p.URL, &p.ImagenURL, &imagenesJSON, // Usamos imagenesJSON en lugar de p.Imagenes directamente
			&p.CreatedAt, &p.UpdatedAt,
			&p.TipoPropiedad, &p.Ubicacion, &p.Dormitorios, &p.Banios, &p.Antiguedad,
			&p.SuperficieCubierta, &p.SuperficieTotal, &p.SuperficieTerreno,
			&p.Frente, &p.Fondo, &p.Ambientes, &p.Plantas, &p.Cocheras,
			&p.Situacion, &p.Expensas, &p.Descripcion, &p.Status, &p.Operacion,
			&p.Condicion, &p.Orientacion, &p.Disposicion,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando propiedad: %v", err)
		}

		// Deserializar JSON si existe
		if imagenesJSON.Valid && imagenesJSON.String != "" {
			var imagenes []string
			if err := json.Unmarshal([]byte(imagenesJSON.String), &imagenes); err != nil {
				return nil, fmt.Errorf("error deserializando imágenes: %v", err)
			}
			p.Imagenes = &imagenes
		}

		propiedades = append(propiedades, p)
	}

	return propiedades, nil
}

// UpdatePropiedadDetalles actualiza solo los campos de detalles de una propiedad
func (db *DB) UpdatePropiedadDetalles(p *Propiedad) error {
	fmt.Printf("Actualizando detalles de propiedad ID: %d\n", p.ID)

	// Convertir el slice de imágenes a JSON
	var imagenesJSON []byte
	var err error
	if p.Imagenes != nil {
		imagenesJSON, err = json.Marshal(*p.Imagenes)
		if err != nil {
			return fmt.Errorf("error convirtiendo imágenes a JSON: %v", err)
		}
	}

	query := `
		UPDATE propiedades 
		SET 
			tipo_propiedad = ?,
			imagenes = ?,
			ubicacion = ?,
			dormitorios = ?,
			banios = ?,
			antiguedad = ?,
			superficie_cubierta = ?,
			superficie_total = ?,
			superficie_terreno = ?,
			frente = ?,
			fondo = ?,
			ambientes = ?,
			plantas = ?,
			cocheras = ?,
			situacion = ?,
			expensas = ?,
			descripcion = ?,
			status = ?,
			operacion = ?,
			condicion = ?,
			orientacion = ?,
			disposicion = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING created_at, updated_at`

	err = db.QueryRow(query,
		p.TipoPropiedad,
		string(imagenesJSON), // Convertimos el JSON a string
		p.Ubicacion,
		p.Dormitorios,
		p.Banios,
		p.Antiguedad,
		p.SuperficieCubierta,
		p.SuperficieTotal,
		p.SuperficieTerreno,
		p.Frente,
		p.Fondo,
		p.Ambientes,
		p.Plantas,
		p.Cocheras,
		p.Situacion,
		p.Expensas,
		p.Descripcion,
		p.Status,
		p.Operacion,
		p.Condicion,
		p.Orientacion,
		p.Disposicion,
		p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error actualizando detalles de propiedad %d: %v", p.ID, err)
	}

	return nil
}

func (db *DB) GetInmobiliariaByID(id int64) (*Inmobiliaria, error) {
	query := `
		SELECT id, nombre, url, sistema, zona, rating, direccion, telefono, created_at, updated_at
		FROM inmobiliarias
		WHERE id = ?`

	var i Inmobiliaria
	err := db.QueryRow(query, id).Scan(
		&i.ID, &i.Nombre, &i.URL, &i.Sistema, &i.Zona, &i.Rating,
		&i.Direccion, &i.Telefono, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo inmobiliaria %d: %v", id, err)
	}

	return &i, nil
}

// GetUnratedProperties retorna las propiedades sin calificar
func (db *DB) GetUnratedProperties(filter *PropertyFilter) ([]Propiedad, error) {
	baseQuery := `
		SELECT 
			p.id, p.inmobiliaria_id, p.codigo, p.titulo, p.precio, p.direccion, 
			p.url, p.imagen_url, p.imagenes, p.created_at, p.updated_at,
			p.tipo_propiedad, p.ubicacion, p.dormitorios, p.banios, p.antiguedad,
			p.superficie_cubierta, p.superficie_total, p.superficie_terreno,
			p.frente, p.fondo, p.ambientes, p.plantas, p.cocheras,
			p.situacion, p.expensas, p.descripcion, p.status, p.operacion,
			p.condicion, p.orientacion, p.disposicion
		FROM propiedades p
		WHERE NOT EXISTS (
			SELECT 1 
			FROM property_ratings r 
			WHERE r.property_id = p.id
		)`

	whereConditions, args := buildFilterConditions(filter)

	// Agregar condiciones WHERE si existen
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	baseQuery += " ORDER BY p.created_at DESC"

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error consultando propiedades sin calificar: %v", err)
	}
	defer rows.Close()

	return scanPropiedades(rows)
}

// GetLikedProperties retorna las propiedades que tienen like
func (db *DB) GetLikedProperties(filter *PropertyFilter) ([]Propiedad, error) {
	baseQuery := `
		SELECT 
			p.id, p.inmobiliaria_id, p.codigo, p.titulo, p.precio, p.direccion, 
			p.url, p.imagen_url, p.imagenes, p.created_at, p.updated_at,
			p.tipo_propiedad, p.ubicacion, p.dormitorios, p.banios, p.antiguedad,
			p.superficie_cubierta, p.superficie_total, p.superficie_terreno,
			p.frente, p.fondo, p.ambientes, p.plantas, p.cocheras,
			p.situacion, p.expensas, p.descripcion, p.status, p.operacion,
			p.condicion, p.orientacion, p.disposicion
		FROM propiedades p
		INNER JOIN property_ratings r ON r.property_id = p.id
		WHERE r.rating = 'like'`

	whereConditions, args := buildFilterConditions(filter)

	// Agregar condiciones WHERE si existen
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	baseQuery += " ORDER BY r.created_at DESC"

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error consultando propiedades con like: %v", err)
	}
	defer rows.Close()

	return scanPropiedades(rows)
}

// Helper para escanear propiedades desde filas de resultados
func scanPropiedades(rows *sql.Rows) ([]Propiedad, error) {
	var propiedades []Propiedad
	for rows.Next() {
		var p Propiedad
		var imagenesJSON sql.NullString // Para manejar NULL en la base de datos

		err := rows.Scan(
			&p.ID, &p.InmobiliariaID, &p.Codigo, &p.Titulo, &p.Precio, &p.Direccion,
			&p.URL, &p.ImagenURL, &imagenesJSON, // Usamos imagenesJSON en lugar de p.Imagenes directamente
			&p.CreatedAt, &p.UpdatedAt,
			&p.TipoPropiedad, &p.Ubicacion, &p.Dormitorios, &p.Banios, &p.Antiguedad,
			&p.SuperficieCubierta, &p.SuperficieTotal, &p.SuperficieTerreno,
			&p.Frente, &p.Fondo, &p.Ambientes, &p.Plantas, &p.Cocheras,
			&p.Situacion, &p.Expensas, &p.Descripcion, &p.Status, &p.Operacion,
			&p.Condicion, &p.Orientacion, &p.Disposicion,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando propiedad: %v", err)
		}

		// Deserializar JSON si existe
		if imagenesJSON.Valid && imagenesJSON.String != "" {
			var imagenes []string
			if err := json.Unmarshal([]byte(imagenesJSON.String), &imagenes); err != nil {
				return nil, fmt.Errorf("error deserializando imágenes: %v", err)
			}
			p.Imagenes = &imagenes
		}

		propiedades = append(propiedades, p)
	}

	return propiedades, nil
}

// buildFilterConditions construye las condiciones WHERE y los argumentos para los filtros
func buildFilterConditions(filter *PropertyFilter) ([]string, []interface{}) {
	if filter == nil {
		return nil, nil
	}

	var conditions []string
	var args []interface{}

	// Filtro por tipo de propiedad
	if filter.PropertyType != "" && filter.PropertyType != "all" {
		conditions = append(conditions, "p.tipo_propiedad = ?")
		args = append(args, filter.PropertyType)
	}

	// Filtro por ubicaciones
	if len(filter.Locations) > 0 {
		placeholders := make([]string, len(filter.Locations))
		for i, loc := range filter.Locations {
			placeholders[i] = "?"
			args = append(args, "%"+loc+"%")
		}
		conditions = append(conditions, fmt.Sprintf("(p.ubicacion LIKE %s)", strings.Join(placeholders, " OR p.ubicacion LIKE ")))
	}

	// Filtro por precio
	if filter.PriceMin != nil || filter.PriceMax != nil {
		// Extraer el valor numérico del precio
		priceExtract := `CAST(REPLACE(REPLACE(REPLACE(p.precio, '$', ''), '.', ''), ',', '') AS NUMERIC)`

		// Ajustar para moneda
		if filter.Currency == "USD" {
			priceExtract = `CASE 
				WHEN p.precio LIKE '%USD%' OR p.precio LIKE '%U$S%' THEN CAST(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(p.precio, 'USD', ''), 'U$S', ''), '$', ''), '.', ''), ',', '') AS NUMERIC)
				ELSE CAST(REPLACE(REPLACE(REPLACE(p.precio, '$', ''), '.', ''), ',', '') AS NUMERIC) / 1000
			END`
		} else {
			// ARS
			priceExtract = `CASE 
				WHEN p.precio LIKE '%USD%' OR p.precio LIKE '%U$S%' THEN CAST(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(p.precio, 'USD', ''), 'U$S', ''), '$', ''), '.', ''), ',', '') AS NUMERIC) * 1000
				ELSE CAST(REPLACE(REPLACE(REPLACE(p.precio, '$', ''), '.', ''), ',', '') AS NUMERIC)
			END`
		}

		if filter.PriceMin != nil {
			conditions = append(conditions, fmt.Sprintf("%s >= ?", priceExtract))
			args = append(args, *filter.PriceMin)
		}

		if filter.PriceMax != nil {
			conditions = append(conditions, fmt.Sprintf("%s <= ?", priceExtract))
			args = append(args, *filter.PriceMax)
		}
	}

	// Filtro por tamaño (superficie)
	if filter.SizeMin != nil {
		conditions = append(conditions, "p.superficie_total >= ?")
		args = append(args, *filter.SizeMin)
	}

	if filter.SizeMax != nil {
		conditions = append(conditions, "p.superficie_total <= ?")
		args = append(args, *filter.SizeMax)
	}

	// Filtro por ambientes
	if filter.Rooms != nil {
		conditions = append(conditions, "p.ambientes >= ?")
		args = append(args, *filter.Rooms)
	}

	// Filtro por baños
	if filter.Bathrooms != nil {
		conditions = append(conditions, "p.banios >= ?")
		args = append(args, *filter.Bathrooms)
	}

	// Filtro por antigüedad
	if filter.Antiquity != nil {
		// Convertir texto de antigüedad a un valor numérico aproximado
		// Esto es una simplificación, idealmente se debería normalizar el campo en la base de datos
		if *filter.Antiquity == 0 {
			conditions = append(conditions, "(p.antiguedad LIKE '%estrenar%' OR p.antiguedad LIKE '%a estrenar%')")
		} else {
			conditions = append(conditions, "p.antiguedad LIKE ?")
			args = append(args, fmt.Sprintf("%%%d%%", *filter.Antiquity))
		}
	}

	// Filtro por características (features)
	// Esto requeriría una tabla de características o buscar en la descripción
	if len(filter.Features) > 0 {
		featureConditions := make([]string, len(filter.Features))
		for i, feature := range filter.Features {
			featureConditions[i] = "p.descripcion LIKE ?"
			args = append(args, "%"+feature+"%")
		}
		conditions = append(conditions, "("+strings.Join(featureConditions, " OR ")+")")
	}

	// Filtro para mostrar solo propiedades con notas
	if filter.ShowOnlyWithNotes {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM property_notes n WHERE n.property_id = p.id
		)`)
	}

	return conditions, args
}

// RateProperty califica una propiedad como like o dislike
func (db *DB) RateProperty(propertyID int64, rating string) error {
	if rating != "like" && rating != "dislike" {
		return fmt.Errorf("rating inválido: %s", rating)
	}

	// Verificar si la propiedad existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM propiedades WHERE id = ?)", propertyID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error verificando propiedad %d: %v", propertyID, err)
	}
	if !exists {
		return fmt.Errorf("la propiedad %d no existe", propertyID)
	}

	// Verificar si ya existe un rating para mantener el estado de favorito
	var isFavorite bool
	var hasRating bool
	err = db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM property_ratings WHERE property_id = ?),
		       IFNULL((SELECT is_favorite FROM property_ratings WHERE property_id = ?), 0)
	`, propertyID, propertyID).Scan(&hasRating, &isFavorite)
	if err != nil {
		return fmt.Errorf("error verificando rating existente: %v", err)
	}

	// Si cambiamos de like a dislike, quitamos el favorito
	if rating == "dislike" {
		isFavorite = false
	}

	// Intentar insertar o actualizar
	query := `
		INSERT INTO property_ratings (property_id, rating, is_favorite)
		VALUES (?, ?, ?)
		ON CONFLICT(property_id) DO UPDATE SET
			rating = excluded.rating,
			is_favorite = CASE 
				WHEN excluded.rating = 'dislike' THEN 0
				ELSE excluded.is_favorite
			END,
			created_at = CURRENT_TIMESTAMP`

	result, err := db.Exec(query, propertyID, rating, isFavorite)
	if err != nil {
		return fmt.Errorf("error calificando propiedad %d: %v", propertyID, err)
	}

	rows, _ := result.RowsAffected()
	fmt.Printf("Filas afectadas: %d\n", rows)

	return nil
}

// GetPropertyNotes obtiene todas las notas de una propiedad
func (db *DB) GetPropertyNotes(propertyID int64) ([]PropertyNote, error) {
	query := `
		SELECT id, property_id, note, created_at, updated_at
		FROM property_notes
		WHERE property_id = ?
		ORDER BY created_at ASC
	`

	// Usar directamente el método Query del sql.DB subyacente
	rows, err := db.DB.Query(query, propertyID)
	if err != nil {
		return nil, fmt.Errorf("error getting property notes: %w", err)
	}
	defer rows.Close()

	notes := []PropertyNote{} // Inicializar como slice vacío en lugar de nil
	for rows.Next() {
		var note PropertyNote
		if err := rows.Scan(&note.ID, &note.PropertyID, &note.Text, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, fmt.Errorf("error scanning property note: %w", err)
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating property notes: %w", err)
	}

	return notes, nil
}

// AddPropertyNote agrega una nota a una propiedad
func (db *DB) AddPropertyNote(note *PropertyNote) error {
	query := `
		INSERT INTO property_notes (property_id, note, created_at, updated_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	result, err := db.DB.Exec(query, note.PropertyID, note.Text)
	if err != nil {
		return fmt.Errorf("error adding property note: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %w", err)
	}

	note.ID = id
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	return nil
}

// DeletePropertyNote elimina una nota de una propiedad
func (db *DB) DeletePropertyNote(noteID int64) error {
	query := `DELETE FROM property_notes WHERE id = ?`

	_, err := db.DB.Exec(query, noteID)
	if err != nil {
		return fmt.Errorf("error deleting property note: %w", err)
	}

	return nil
}

// PropertyHasNotes verifica si una propiedad tiene notas
func (db *DB) PropertyHasNotes(propertyID int64) (bool, error) {
	query := `SELECT COUNT(*) FROM property_notes WHERE property_id = ?`

	var count int
	err := db.DB.QueryRow(query, propertyID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if property has notes: %w", err)
	}

	return count > 0, nil
}

// TogglePropertyFavorite marca o desmarca una propiedad como favorita
func (db *DB) TogglePropertyFavorite(propertyID int64, isFavorite bool) error {
	// Verificar si la propiedad existe
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM propiedades WHERE id = ?)", propertyID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("error verificando propiedad %d: %v", propertyID, err)
	}
	if !exists {
		return fmt.Errorf("la propiedad %d no existe", propertyID)
	}

	// Verificar si la propiedad ya tiene un rating
	var hasRating bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM property_ratings WHERE property_id = ?)", propertyID).Scan(&hasRating)
	if err != nil {
		return fmt.Errorf("error verificando rating de propiedad %d: %v", propertyID, err)
	}

	if hasRating {
		// Actualizar el valor de is_favorite
		query := `
			UPDATE property_ratings
			SET is_favorite = ?
			WHERE property_id = ?`

		_, err = db.Exec(query, isFavorite, propertyID)
		if err != nil {
			return fmt.Errorf("error actualizando favorito para propiedad %d: %v", propertyID, err)
		}
	} else {
		// Si no tiene rating, no podemos marcarla como favorita (debe tener like primero)
		return fmt.Errorf("la propiedad %d no tiene calificación, debe tener 'like' antes de marcarla como favorita", propertyID)
	}

	return nil
}

// GetFavoriteProperties retorna las propiedades marcadas como favoritas
func (db *DB) GetFavoriteProperties(filter *PropertyFilter) ([]Propiedad, error) {
	baseQuery := `
		SELECT 
			p.id, p.inmobiliaria_id, p.codigo, p.titulo, p.precio, p.direccion, 
			p.url, p.imagen_url, p.imagenes, p.created_at, p.updated_at,
			p.tipo_propiedad, p.ubicacion, p.dormitorios, p.banios, p.antiguedad,
			p.superficie_cubierta, p.superficie_total, p.superficie_terreno,
			p.frente, p.fondo, p.ambientes, p.plantas, p.cocheras,
			p.situacion, p.expensas, p.descripcion, p.status, p.operacion,
			p.condicion, p.orientacion, p.disposicion
		FROM propiedades p
		INNER JOIN property_ratings r ON r.property_id = p.id
		WHERE r.rating = 'like' AND r.is_favorite = 1`

	whereConditions, args := buildFilterConditions(filter)

	// Agregar condiciones WHERE si existen
	if len(whereConditions) > 0 {
		baseQuery += " AND " + strings.Join(whereConditions, " AND ")
	}

	baseQuery += " ORDER BY r.created_at DESC"

	rows, err := db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("error consultando propiedades favoritas: %v", err)
	}
	defer rows.Close()

	return scanPropiedades(rows)
}

// IsPropertyFavorite verifica si una propiedad está marcada como favorita
func (db *DB) IsPropertyFavorite(propertyID int64) (bool, error) {
	query := `SELECT is_favorite FROM property_ratings WHERE property_id = ?`

	var isFavorite bool
	err := db.QueryRow(query, propertyID).Scan(&isFavorite)
	if err == sql.ErrNoRows {
		// Si no hay registro, no es favorita
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("error verificando si la propiedad %d es favorita: %v", propertyID, err)
	}

	return isFavorite, nil
}
