package db

import (
	"database/sql"
	"embed"
	"fmt"
	"strings"

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

	return &DB{db}, nil
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
            inmobiliaria_id, codigo, titulo, precio, direccion, url, imagen_url, fecha_scraping,
            tipo_propiedad, ubicacion, dormitorios, banios, antiguedad, 
            superficie_cubierta, superficie_total, frente, fondo, ambientes,
            expensas, descripcion, status
        )
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(codigo) DO UPDATE SET
            titulo = excluded.titulo,
            precio = excluded.precio,
            direccion = excluded.direccion,
            url = excluded.url,
            imagen_url = excluded.imagen_url,
            fecha_scraping = excluded.fecha_scraping,
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
		p.InmobiliariaID, p.Codigo, p.Titulo, p.Precio, p.Direccion, p.URL, p.ImagenURL, p.FechaScraping,
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

	// Debug logs...

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
			fecha_scraping, created_at, updated_at,
			tipo_propiedad, ubicacion, dormitorios, banios, antiguedad,
			superficie_cubierta, superficie_total, frente, fondo, ambientes,
			expensas, descripcion, status
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
		err := rows.Scan(
			&p.ID, &p.InmobiliariaID, &p.Codigo, &p.Titulo, &p.Precio, &p.Direccion,
			&p.URL, &p.ImagenURL, &p.FechaScraping, &p.CreatedAt, &p.UpdatedAt,
			&p.TipoPropiedad, &p.Ubicacion, &p.Dormitorios, &p.Banios, &p.Antiguedad,
			&p.SuperficieCubierta, &p.SuperficieTotal, &p.Frente, &p.Fondo, &p.Ambientes,
			&p.Expensas, &p.Descripcion, &p.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando propiedad: %v", err)
		}
		propiedades = append(propiedades, p)
	}

	return propiedades, nil
}

// UpdatePropiedadDetalles actualiza solo los campos de detalles de una propiedad
func (db *DB) UpdatePropiedadDetalles(p *Propiedad) error {
	query := `
		UPDATE propiedades 
		SET 
			tipo_propiedad = ?,
			ubicacion = ?,
			dormitorios = ?,
			banios = ?,
			antiguedad = ?,
			superficie_cubierta = ?,
			superficie_total = ?,
			frente = ?,
			fondo = ?,
			ambientes = ?,
			expensas = ?,
			descripcion = ?,
			status = ?,
			updated_at = CURRENT_TIMESTAMP,
			fecha_scraping = ?
		WHERE id = ?
		RETURNING created_at, updated_at`

	err := db.QueryRow(query,
		p.TipoPropiedad,
		p.Ubicacion,
		p.Dormitorios,
		p.Banios,
		p.Antiguedad,
		p.SuperficieCubierta,
		p.SuperficieTotal,
		p.Frente,
		p.Fondo,
		p.Ambientes,
		p.Expensas,
		p.Descripcion,
		p.Status,
		p.FechaScraping,
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
