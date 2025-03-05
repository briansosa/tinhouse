package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/findhouse/internal/db"
)

// normalizarTipoPropiedad normaliza el tipo de propiedad para que coincida con los códigos en la base de datos
func normalizarTipoPropiedad(tipo string) string {
	if tipo == "" {
		return ""
	}

	// Convertir a minúsculas y eliminar espacios
	tipo = strings.TrimSpace(tipo)

	// Mapeo de tipos comunes
	switch strings.ToLower(tipo) {
	case "casa", "chalet", "casa chalet", "casa quinta", "quinta":
		return "Casa"
	case "departamento", "depto", "dpto", "dpto.", "depto.", "departamento con dependencia":
		return "Departamento"
	case "ph", "p.h.", "p.h", "propiedad horizontal":
		return "PH"
	case "local", "local comercial", "fondo de comercio":
		return "Local"
	case "oficina", "consultorio", "estudio":
		return "Oficina"
	case "terreno", "lote", "lote de terreno", "fracción", "fraccion":
		return "Terreno"
	case "galpón", "galpon", "depósito", "deposito", "nave industrial":
		return "Galpón"
	default:
		// Si no hay coincidencia, devolvemos el tipo original
		return tipo
	}
}

func main() {
	dbPath := flag.String("db", "../../internal/db/findhouse.db", "Path to SQLite database")
	flag.Parse()

	// Inicializar DB
	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer database.Close()

	// Crear tabla de tipos de propiedad si no existe
	_, err = database.Exec(`
		CREATE TABLE IF NOT EXISTS property_types (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal("Error creating property_types table:", err)
	}

	// Insertar tipos de propiedad predeterminados
	_, err = database.Exec(`
		INSERT OR IGNORE INTO property_types (code, name) VALUES 
		('house', 'Casa'),
		('apartment', 'Departamento'),
		('ph', 'PH'),
		('local', 'Local'),
		('office', 'Oficina'),
		('land', 'Terreno'),
		('warehouse', 'Galpón')
	`)
	if err != nil {
		log.Fatal("Error inserting default property types:", err)
	}

	// Obtener todos los tipos de propiedad distintos de la base de datos
	rows, err := database.Query(`
		SELECT DISTINCT tipo_propiedad 
		FROM propiedades 
		WHERE tipo_propiedad IS NOT NULL
	`)
	if err != nil {
		log.Fatal("Error querying property types:", err)
	}
	defer rows.Close()

	// Mapear y normalizar los tipos de propiedad
	tiposEncontrados := make(map[string]string) // original -> normalizado
	for rows.Next() {
		var tipoOriginal sql.NullString
		if err := rows.Scan(&tipoOriginal); err != nil {
			log.Fatal("Error scanning property type:", err)
		}

		if tipoOriginal.Valid && tipoOriginal.String != "" {
			tipoNormalizado := normalizarTipoPropiedad(tipoOriginal.String)
			tiposEncontrados[tipoOriginal.String] = tipoNormalizado

			// Insertar en la tabla de tipos si no existe
			var code string
			switch tipoNormalizado {
			case "Casa":
				code = "house"
			case "Departamento":
				code = "apartment"
			case "PH":
				code = "ph"
			case "Local":
				code = "local"
			case "Oficina":
				code = "office"
			case "Terreno":
				code = "land"
			case "Galpón":
				code = "warehouse"
			default:
				// Generar un código a partir del nombre normalizado
				code = strings.ToLower(strings.ReplaceAll(tipoNormalizado, " ", "_"))
			}

			_, err = database.Exec(`
				INSERT OR IGNORE INTO property_types (code, name) 
				VALUES (?, ?)
			`, code, tipoNormalizado)
			if err != nil {
				log.Printf("Error inserting property type %s: %v", tipoNormalizado, err)
			}
		}
	}

	// Mostrar los tipos encontrados y sus normalizaciones
	fmt.Println("Tipos de propiedad encontrados y normalizados:")
	for original, normalizado := range tiposEncontrados {
		fmt.Printf("  %s -> %s\n", original, normalizado)
	}

	// Actualizar las propiedades con los tipos normalizados
	for original, normalizado := range tiposEncontrados {
		_, err = database.Exec(`
			UPDATE propiedades 
			SET tipo_propiedad = ? 
			WHERE tipo_propiedad = ?
		`, normalizado, original)
		if err != nil {
			log.Printf("Error updating properties with type %s: %v", original, err)
		} else {
			fmt.Printf("Actualizado: %s -> %s\n", original, normalizado)
		}
	}

	fmt.Println("Migración completada con éxito.")
}
