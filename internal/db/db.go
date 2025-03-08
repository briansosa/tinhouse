package db

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
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
			descripcion, status, operacion, condicion, orientacion, disposicion,
			latitud, longitud
		FROM propiedades
		WHERE status = 'pending'
		ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando propiedades sin detalles: %v", err)
	}
	defer rows.Close()

	return scanPropiedades(rows)
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

	// Comenzar una transacción
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error iniciando transacción: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
	}()

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
			latitud = ?,
			longitud = ?,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
		RETURNING created_at, updated_at`

	err = tx.QueryRow(query,
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
		p.Latitud,
		p.Longitud,
		p.ID,
	).Scan(&p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return fmt.Errorf("error actualizando detalles de propiedad %d: %v", p.ID, err)
	}

	// Si hay características para guardar, las guardamos usando la misma transacción
	if p.Features != nil && len(p.Features) > 0 {
		if err := db.SavePropertyFeaturesWithTx(tx, p.ID, p.Features); err != nil {
			return fmt.Errorf("error guardando características: %v", err)
		}
	}

	// Confirmar la transacción
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error confirmando transacción: %v", err)
	}

	return nil
}

// SavePropertyFeaturesWithTx guarda las características de una propiedad usando una transacción existente
func (db *DB) SavePropertyFeaturesWithTx(tx *sql.Tx, propertyID int64, features map[string][]string) error {
	// Eliminar relaciones existentes dentro de la transacción
	deleteQuery := `DELETE FROM property_feature_relations WHERE property_id = ?`
	_, err := tx.Exec(deleteQuery, propertyID)
	if err != nil {
		return fmt.Errorf("error eliminando características existentes: %v", err)
	}

	// Para cada categoría y característica, buscamos su ID o la creamos si no existe
	for category, names := range features {
		for _, name := range names {
			// Normalizar el nombre (trim y capitalizar primera letra)
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}

			// Buscar la característica por nombre y categoría
			var featureID int64
			query := `SELECT id FROM property_features WHERE name = ? AND category = ?`
			err := tx.QueryRow(query, name, category).Scan(&featureID)

			if err == sql.ErrNoRows {
				// La característica no existe, la creamos usando INSERT OR IGNORE para evitar errores de duplicados
				insertQuery := `INSERT OR IGNORE INTO property_features (name, category) VALUES (?, ?)`
				result, err := tx.Exec(insertQuery, name, category)
				if err != nil {
					// Si hay un error que no sea de duplicado, lo reportamos
					if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
						return fmt.Errorf("error creando característica %s (%s): %v", name, category, err)
					}

					// Si fue un error de duplicado, intentamos obtener el ID existente
					err = tx.QueryRow(query, name, category).Scan(&featureID)
					if err != nil {
						return fmt.Errorf("error recuperando ID después de duplicado para %s (%s): %v", name, category, err)
					}
				} else {
					// Si no hubo error, obtenemos el ID de la inserción
					newID, err := result.LastInsertId()
					if err != nil {
						return fmt.Errorf("error obteniendo ID de característica: %v", err)
					}

					// Solo asignamos el nuevo ID si realmente se insertó una fila
					if newID > 0 {
						featureID = newID
					} else {
						// Si no se insertó (porque ya existía), obtenemos el ID existente
						err = tx.QueryRow(query, name, category).Scan(&featureID)
						if err != nil {
							return fmt.Errorf("error recuperando ID existente para %s (%s): %v", name, category, err)
						}
					}
				}
			} else if err != nil {
				return fmt.Errorf("error buscando característica %s (%s): %v", name, category, err)
			}

			// Crear la relación entre la propiedad y la característica
			relationQuery := `INSERT OR IGNORE INTO property_feature_relations (property_id, feature_id) VALUES (?, ?)`
			_, err = tx.Exec(relationQuery, propertyID, featureID)
			if err != nil {
				return fmt.Errorf("error creando relación para característica %d: %v", featureID, err)
			}
		}
	}

	return nil
}

// SavePropertyFeatures guarda las características de una propiedad
func (db *DB) SavePropertyFeatures(propertyID int64, features map[string][]string) error {
	// Implementar reintentos con backoff exponencial
	maxRetries := 5
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Si no es el primer intento, esperar con backoff exponencial
		if attempt > 0 {
			backoffTime := time.Duration(math.Pow(2, float64(attempt))) * 100 * time.Millisecond
			time.Sleep(backoffTime)
			fmt.Printf("Reintentando guardar características (intento %d/%d) después de %v\n",
				attempt+1, maxRetries, backoffTime)
		}

		// Iniciar una transacción para asegurar consistencia
		tx, err := db.Begin()
		if err != nil {
			lastErr = fmt.Errorf("error iniciando transacción: %v", err)
			continue
		}

		// Usar la función con transacción
		err = db.SavePropertyFeaturesWithTx(tx, propertyID, features)
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), "database is locked") {
				lastErr = fmt.Errorf("base de datos bloqueada: %v", err)
				continue // Reintentar
			}
			return err
		}

		// Confirmar la transacción
		if err := tx.Commit(); err != nil {
			if strings.Contains(err.Error(), "database is locked") {
				lastErr = fmt.Errorf("base de datos bloqueada al confirmar transacción: %v", err)
				continue // Reintentar
			}
			return fmt.Errorf("error confirmando transacción: %v", err)
		}

		// Si llegamos aquí, todo salió bien
		return nil
	}

	// Si llegamos aquí, agotamos los reintentos
	return fmt.Errorf("error guardando características después de %d intentos: %v", maxRetries, lastErr)
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
			p.condicion, p.orientacion, p.disposicion, p.latitud, p.longitud
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

	propiedades, err := scanPropiedades(rows)
	return propiedades, err
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
			p.condicion, p.orientacion, p.disposicion, p.latitud, p.longitud
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
	var count int

	for rows.Next() {
		count++
		var p Propiedad
		var imagenes sql.NullString

		err := rows.Scan(
			&p.ID, &p.InmobiliariaID, &p.Codigo, &p.Titulo, &p.Precio, &p.Direccion,
			&p.URL, &p.ImagenURL, &imagenes, &p.CreatedAt, &p.UpdatedAt,
			&p.TipoPropiedad, &p.Ubicacion, &p.Dormitorios, &p.Banios, &p.Antiguedad,
			&p.SuperficieCubierta, &p.SuperficieTotal, &p.SuperficieTerreno,
			&p.Frente, &p.Fondo, &p.Ambientes, &p.Plantas, &p.Cocheras,
			&p.Situacion, &p.Expensas, &p.Descripcion, &p.Status, &p.Operacion,
			&p.Condicion, &p.Orientacion, &p.Disposicion, &p.Latitud, &p.Longitud,
		)

		if err != nil {
			fmt.Printf("Error escaneando propiedad #%d: %v\n", count, err)
			return nil, fmt.Errorf("error escaneando propiedad: %v", err)
		}

		// Convertir string de imágenes a slice
		if imagenes.Valid {
			var imgs []string
			if err := json.Unmarshal([]byte(imagenes.String), &imgs); err != nil {
				fmt.Printf("Error deserializando imágenes para propiedad #%d: %v\n", count, err)
				p.Imagenes = &[]string{}
			} else {
				p.Imagenes = &imgs
			}
		} else {
			p.Imagenes = &[]string{}
		}

		propiedades = append(propiedades, p)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error iterando filas: %v\n", err)
		return nil, fmt.Errorf("error iterando propiedades: %v", err)
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
	if filter.PropertyTypeIDs != nil && len(filter.PropertyTypeIDs) > 0 {
		// Usar múltiples IDs de tipos de propiedad
		placeholders := make([]string, len(filter.PropertyTypeIDs))
		for i := range filter.PropertyTypeIDs {
			placeholders[i] = "?"
			args = append(args, filter.PropertyTypeIDs[i])
		}
		conditions = append(conditions, fmt.Sprintf("p.tipo_propiedad IN (%s)", strings.Join(placeholders, ",")))
		fmt.Printf("Agregando condición de tipo_propiedad con múltiples IDs: %v\n", filter.PropertyTypeIDs)
	} else if filter.PropertyTypeID != nil {
		// Usar directamente el ID del tipo de propiedad
		conditions = append(conditions, "p.tipo_propiedad = ?")
		args = append(args, *filter.PropertyTypeID)
		fmt.Printf("Agregando condición de tipo_propiedad con ID: %d\n", *filter.PropertyTypeID)
	} else if filter.PropertyType != "" && filter.PropertyType != "all" {
		// Mantener compatibilidad con el filtro por código
		conditions = append(conditions, "p.tipo_propiedad = (SELECT id FROM property_types WHERE code = ?)")
		args = append(args, filter.PropertyType)
		fmt.Printf("Agregando condición de tipo_propiedad con código: %s\n", filter.PropertyType)
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
		// Primero filtramos por la moneda seleccionada
		if filter.Currency == "USD" {
			conditions = append(conditions, "(p.precio LIKE '%USD%' OR p.precio LIKE '%U$S%' OR p.moneda = 'USD')")
		} else {
			// Si es ARS, excluimos las propiedades en USD
			conditions = append(conditions, "(p.precio NOT LIKE '%USD%' AND p.precio NOT LIKE '%U$S%' AND (p.moneda IS NULL OR p.moneda = 'ARS'))")
		}

		// Extraer el valor numérico del precio
		priceExtract := `CAST(REPLACE(REPLACE(REPLACE(REPLACE(REPLACE(p.precio, 'USD', ''), 'U$S', ''), '$', ''), '.', ''), ',', '') AS NUMERIC)`

		if filter.PriceMin != nil {
			conditions = append(conditions, fmt.Sprintf("%s >= ?", priceExtract))
			args = append(args, *filter.PriceMin)
		}

		if filter.PriceMax != nil {
			conditions = append(conditions, fmt.Sprintf("%s <= ?", priceExtract))
			args = append(args, *filter.PriceMax)
		}
	} else if filter.Currency != "" {
		// Si solo se especifica la moneda sin rango de precios
		if filter.Currency == "USD" {
			conditions = append(conditions, "(p.precio LIKE '%USD%' OR p.precio LIKE '%U$S%' OR p.moneda = 'USD')")
		} else {
			conditions = append(conditions, "(p.precio NOT LIKE '%USD%' AND p.precio NOT LIKE '%U$S%' AND (p.moneda IS NULL OR p.moneda = 'ARS'))")
		}
	}

	// Filtro por tamaño (superficie) - Compatibilidad con versión anterior
	if filter.SizeMin != nil {
		conditions = append(conditions, "(p.superficie_total >= ? OR p.superficie_cubierta >= ?)")
		args = append(args, *filter.SizeMin, *filter.SizeMin)
	}

	if filter.SizeMax != nil {
		conditions = append(conditions, "(p.superficie_total <= ? OR p.superficie_cubierta <= ?)")
		args = append(args, *filter.SizeMax, *filter.SizeMax)
	}

	// Filtros de superficie específicos

	// Superficie Total
	if filter.TotalAreaMin != nil {
		conditions = append(conditions, "p.superficie_total >= ?")
		args = append(args, *filter.TotalAreaMin)
	}

	if filter.TotalAreaMax != nil {
		conditions = append(conditions, "p.superficie_total <= ?")
		args = append(args, *filter.TotalAreaMax)
	}

	// Superficie Cubierta
	if filter.CoveredAreaMin != nil {
		conditions = append(conditions, "p.superficie_cubierta >= ?")
		args = append(args, *filter.CoveredAreaMin)
	}

	if filter.CoveredAreaMax != nil {
		conditions = append(conditions, "p.superficie_cubierta <= ?")
		args = append(args, *filter.CoveredAreaMax)
	}

	// Superficie Terreno
	if filter.LandAreaMin != nil {
		conditions = append(conditions, "p.superficie_terreno >= ?")
		args = append(args, *filter.LandAreaMin)
	}

	if filter.LandAreaMax != nil {
		conditions = append(conditions, "p.superficie_terreno <= ?")
		args = append(args, *filter.LandAreaMax)
	}

	// Frente
	if filter.Front != nil {
		conditions = append(conditions, "p.frente >= ?")
		args = append(args, *filter.Front)
	}

	// Fondo
	if filter.Back != nil {
		conditions = append(conditions, "p.fondo >= ?")
		args = append(args, *filter.Back)
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
		// Ahora que antiguedad es un campo numérico, podemos usar operadores de comparación directos
		switch *filter.Antiquity {
		case 0:
			// A estrenar
			conditions = append(conditions, "p.antiguedad = 0")
		case 5:
			// Hasta 5 años
			conditions = append(conditions, "p.antiguedad > 0 AND p.antiguedad <= 5")
		case 10:
			// Hasta 10 años
			conditions = append(conditions, "p.antiguedad > 0 AND p.antiguedad <= 10")
		case 20:
			// Hasta 20 años
			conditions = append(conditions, "p.antiguedad > 0 AND p.antiguedad <= 20")
		case 30:
			// Hasta 30 años
			conditions = append(conditions, "p.antiguedad > 0 AND p.antiguedad <= 30")
		case 100:
			// Más de 30 años
			conditions = append(conditions, "p.antiguedad > 30")
		default:
			// Caso genérico
			conditions = append(conditions, "p.antiguedad = ?")
			args = append(args, *filter.Antiquity)
		}
	}

	// Filtro por características (features)
	if len(filter.Features) > 0 {
		fmt.Printf("Procesando filtro de características: %v\n", filter.Features)

		// Usar la tabla de relaciones de características con IDs numéricos
		featureCondition := `p.id IN (
			SELECT DISTINCT pf.property_id 
			FROM property_feature_relations pf 
			WHERE pf.feature_id IN (`

		placeholders := make([]string, len(filter.Features))
		featureIDs := make([]int64, 0, len(filter.Features))

		for i := range filter.Features {
			placeholders[i] = "?"
			// Convertir string a int64
			featureID, err := strconv.ParseInt(filter.Features[i], 10, 64)
			if err != nil {
				// Si no es un número válido, usar 0 (que no coincidirá con ningún ID)
				featureID = 0
				fmt.Printf("Error convirtiendo feature ID '%s': %v\n", filter.Features[i], err)
			}
			featureIDs = append(featureIDs, featureID)
			args = append(args, featureID)
		}

		fmt.Printf("Filtrando por features IDs: %v\n", featureIDs)

		featureCondition += strings.Join(placeholders, ", ") + "))"
		conditions = append(conditions, featureCondition)
	}

	// Filtro para mostrar solo propiedades con notas
	if filter.ShowOnlyWithNotes {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM property_notes n WHERE n.property_id = p.id
		)`)
	}

	// Filtro para mostrar solo propiedades favoritas
	if filter.ShowOnlyFavorites {
		conditions = append(conditions, `EXISTS (
			SELECT 1 FROM property_ratings r WHERE r.property_id = p.id AND r.is_favorite = 1
		)`)
	}

	// Filtro por disposición
	if filter.Disposition != nil && len(filter.Disposition) > 0 {
		placeholders := make([]string, len(filter.Disposition))
		for i, disp := range filter.Disposition {
			placeholders[i] = "?"
			args = append(args, disp)
		}
		conditions = append(conditions, fmt.Sprintf("p.disposicion IN (%s)", strings.Join(placeholders, ",")))
	}

	// Filtro por orientación
	if filter.Orientation != nil && len(filter.Orientation) > 0 {
		placeholders := make([]string, len(filter.Orientation))
		for i, orient := range filter.Orientation {
			placeholders[i] = "?"
			args = append(args, orient)
		}
		conditions = append(conditions, fmt.Sprintf("p.orientacion IN (%s)", strings.Join(placeholders, ",")))
	}

	// Filtro por condición
	if filter.Condition != nil && len(filter.Condition) > 0 {
		placeholders := make([]string, len(filter.Condition))
		for i, cond := range filter.Condition {
			placeholders[i] = "?"
			args = append(args, cond)
		}
		conditions = append(conditions, fmt.Sprintf("p.condicion IN (%s)", strings.Join(placeholders, ",")))
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
			p.condicion, p.orientacion, p.disposicion, p.latitud, p.longitud
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

// GetPropertyFeatures obtiene todas las características de una propiedad
func (db *DB) GetPropertyFeatures(propertyID int64) (map[string][]PropertyFeature, error) {
	query := `
		SELECT f.id, f.name, f.category, f.created_at
		FROM property_features f
		JOIN property_feature_relations r ON f.id = r.feature_id
		WHERE r.property_id = ?
		ORDER BY f.category, f.name
	`

	rows, err := db.Query(query, propertyID)
	if err != nil {
		return nil, fmt.Errorf("error obteniendo características de la propiedad %d: %v", propertyID, err)
	}
	defer rows.Close()

	// Agrupar características por categoría
	result := make(map[string][]PropertyFeature)
	for rows.Next() {
		var feature PropertyFeature
		if err := rows.Scan(&feature.ID, &feature.Name, &feature.Category, &feature.CreatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando característica: %v", err)
		}

		result[feature.Category] = append(result[feature.Category], feature)
	}

	return result, nil
}

// GetPropertyFeaturesAsMap obtiene todas las características de una propiedad como un mapa
func (db *DB) GetPropertyFeaturesAsMap(propertyID int64) (map[string][]string, error) {
	features, err := db.GetPropertyFeatures(propertyID)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for category, featureList := range features {
		for _, feature := range featureList {
			result[category] = append(result[category], feature.Name)
		}
	}

	return result, nil
}

// GetAllFeatures retorna todas las características disponibles
func (db *DB) GetAllFeatures() ([]PropertyFeature, error) {
	query := `SELECT id, name, category, created_at FROM property_features ORDER BY category, name`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando características: %v", err)
	}
	defer rows.Close()

	var features []PropertyFeature
	for rows.Next() {
		var feature PropertyFeature
		if err := rows.Scan(&feature.ID, &feature.Name, &feature.Category, &feature.CreatedAt); err != nil {
			return nil, fmt.Errorf("error escaneando característica: %v", err)
		}
		features = append(features, feature)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando características: %v", err)
	}

	return features, nil
}

// GetAllPropertyTypes obtiene todos los tipos de propiedad disponibles
func (db *DB) GetAllPropertyTypes() ([]PropertyType, error) {
	var types []PropertyType
	query := `SELECT id, code, name, created_at FROM property_types ORDER BY name`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Printf("Error en la consulta: %v\n", err)
		return nil, fmt.Errorf("error al obtener tipos de propiedad: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pt PropertyType
		if err := rows.Scan(&pt.ID, &pt.Code, &pt.Name, &pt.CreatedAt); err != nil {
			fmt.Printf("Error al escanear fila: %v\n", err)
			return nil, fmt.Errorf("error al escanear tipo de propiedad: %w", err)
		}
		types = append(types, pt)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("Error al iterar filas: %v\n", err)
		return nil, fmt.Errorf("error al iterar tipos de propiedad: %w", err)
	}

	return types, nil
}

// GetPropertyTypeByCode obtiene un tipo de propiedad por su código
func (db *DB) GetPropertyTypeByCode(code string) (*PropertyType, error) {
	var propertyType PropertyType
	query := `SELECT id, code, name, created_at FROM property_types WHERE code = ?`

	row := db.QueryRow(query, code)
	err := row.Scan(&propertyType.ID, &propertyType.Code, &propertyType.Name, &propertyType.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No se encontró, pero no es un error
		}
		return nil, fmt.Errorf("error al obtener tipo de propiedad por código: %w", err)
	}

	return &propertyType, nil
}

// GetPropertyTypeNameByCode obtiene el nombre de un tipo de propiedad por su código
func (db *DB) GetPropertyTypeNameByCode(code string) (string, error) {
	if code == "" || code == "all" {
		return "", nil
	}

	var name string
	query := `SELECT name FROM property_types WHERE code = ?`

	row := db.QueryRow(query, code)
	err := row.Scan(&name)
	if err != nil {
		if err == sql.ErrNoRows {
			return code, nil // Si no existe, devolvemos el código como fallback
		}
		return "", fmt.Errorf("error al obtener nombre de tipo de propiedad: %w", err)
	}

	return name, nil
}

// GetListValuesByName devuelve los valores de una lista específica
func (db *DB) GetListValuesByName(listName string) ([]ListValue, error) {
	query := `
		SELECT lv.id, lv.list_id, lv.value, lv.display_name, lv.sort_order, lv.created_at, lv.updated_at
		FROM list_values lv
		JOIN lists l ON lv.list_id = l.id
		WHERE l.name = ?
		ORDER BY lv.sort_order, lv.display_name
	`

	rows, err := db.Query(query, listName)
	if err != nil {
		return nil, fmt.Errorf("error consultando valores de lista %s: %w", listName, err)
	}
	defer rows.Close()

	var values []ListValue
	for rows.Next() {
		var value ListValue
		err := rows.Scan(
			&value.ID,
			&value.ListID,
			&value.Value,
			&value.DisplayName,
			&value.SortOrder,
			&value.CreatedAt,
			&value.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando valor de lista: %w", err)
		}
		values = append(values, value)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %w", err)
	}

	// Si no hay valores en la base de datos, devolver valores por defecto
	if len(values) == 0 {
		switch listName {
		case "disposition":
			return createDefaultListValues(listName, []string{"frente", "contrafrente", "interno", "lateral"}), nil
		case "orientation":
			return createDefaultListValues(listName, []string{"norte", "sur", "este", "oeste", "noreste", "noroeste", "sureste", "suroeste"}), nil
		case "status":
			return createDefaultListValues(listName, []string{"a estrenar", "a reciclar", "en construcción", "refaccionado", "excelente"}), nil
		case "operation":
			return createDefaultListValues(listName, []string{"venta", "alquiler", "alquiler temporario"}), nil
		case "condition":
			return createDefaultListValues(listName, []string{"vacía", "ocupada", "en sucesión"}), nil
		}
	}

	return values, nil
}

// createDefaultListValues crea valores de lista por defecto cuando no existen en la base de datos
func createDefaultListValues(listName string, values []string) []ListValue {
	result := make([]ListValue, len(values))
	for i, v := range values {
		result[i] = ListValue{
			ID:          int64(i + 1),
			Value:       v,
			DisplayName: strings.Title(v),
			SortOrder:   i + 1,
		}
	}
	return result
}

// PropertyFeatureRelation representa la relación entre una propiedad y una característica
