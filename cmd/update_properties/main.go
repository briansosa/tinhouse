package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/findhouse/internal/db"
	"github.com/findhouse/internal/scraper/tokko"
)

func main() {
	dbPath := flag.String("db", "findhouse.db", "Ruta a la base de datos SQLite")
	testMode := flag.Bool("test", false, "Ejecutar en modo test (solo una propiedad)")
	flag.Parse()

	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer database.Close()

	if err := updateProperties(database, *testMode); err != nil {
		log.Fatal(err)
	}
}

func updateProperties(database *db.DB, testMode bool) error {
	fmt.Println("🚀 Iniciando proceso de actualización")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Println("✓ Contexto creado con timeout de 5 minutos")

	// Obtener propiedades pendientes
	propiedades, err := database.GetPropiedadesSinDetalles()
	if err != nil {
		return fmt.Errorf("error obteniendo propiedades sin detalles: %v", err)
	}

	fmt.Printf("Encontradas %d propiedades pendientes\n", len(propiedades))

	if testMode && len(propiedades) > 0 {
		propiedades = propiedades[:1]
		fmt.Println("Modo test: usando solo una propiedad")
	}

	var actualizadas, fallidas, noDisponibles int

	for _, prop := range propiedades {
		fmt.Printf("\n📝 Procesando propiedad %s\n", prop.Codigo)
		fmt.Printf("   URL: %s\n", prop.URL)
		fmt.Printf("   Inmobiliaria ID: %d\n", prop.InmobiliariaID)

		inmobiliaria, err := database.GetInmobiliariaByID(prop.InmobiliariaID)
		if err != nil {
			log.Printf("❌ Error obteniendo inmobiliaria %d: %v\n", prop.InmobiliariaID, err)
			fallidas++
			continue
		}

		fmt.Printf("   Inmobiliaria: %s (Sistema: %s)\n", inmobiliaria.Nombre, inmobiliaria.Sistema)

		// Por ahora solo manejamos Tokko
		if !strings.Contains(strings.ToLower(inmobiliaria.Sistema), "tokko") {
			log.Printf("Sistema no soportado: %s\n", inmobiliaria.Sistema)
			fallidas++
			continue
		}

		// Obtener detalles usando el contexto con timeout
		scraper := tokko.New(inmobiliaria.URL)
		details, err := scraper.GetPropertyDetails(ctx, prop.URL)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				log.Printf("Propiedad no disponible: %s\n", prop.URL)
				prop.Status = "not_available"
				noDisponibles++
			} else {
				log.Printf("Error obteniendo detalles: %v\n", err)
				fallidas++
				continue
			}
		} else {
			// Actualizar propiedad con los detalles
			prop.TipoPropiedad = &details.TipoPropiedad
			prop.Ubicacion = &details.Ubicacion
			prop.Dormitorios = &details.Dormitorios
			prop.Banios = &details.Banios
			prop.Antiguedad = &details.Antiguedad
			prop.SuperficieCubierta = &details.SuperficieCubierta
			prop.Frente = &details.Frente
			prop.Fondo = &details.Fondo
			prop.Ambientes = &details.Ambientes
			prop.Expensas = &details.Expensas
			prop.Descripcion = &details.Descripcion
			prop.Status = "match"
			actualizadas++
		}

		prop.FechaScraping = time.Now()
		if err := database.UpdatePropiedadDetalles(&prop); err != nil {
			log.Printf("Error actualizando propiedad: %v\n", err)
			fallidas++
			continue
		}

		time.Sleep(1 * time.Second) // Delay entre requests
	}

	fmt.Printf("\nResumen:\n"+
		"- Propiedades actualizadas: %d\n"+
		"- Propiedades fallidas: %d\n"+
		"- Propiedades no disponibles: %d\n",
		actualizadas, fallidas, noDisponibles)

	return nil
}
