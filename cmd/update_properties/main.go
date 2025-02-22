package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
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

// Función auxiliar para convertir a puntero
func ptr[T any](v T) *T {
	return &v
}

func updateProperties(database *db.DB, testMode bool) error {
	fmt.Println("🚀 Iniciando proceso de actualización")

	// Quitamos el timeout global
	ctx := context.Background()

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
		propCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		details, err := scraper.GetPropertyDetails(propCtx, prop.URL)
		cancel() // Liberamos recursos

		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				log.Printf("⌛ Timeout al procesar propiedad: %s, reintentando...\n", prop.URL)
				// Reintento con más tiempo
				propCtx, cancel = context.WithTimeout(ctx, 5*time.Minute)
				details, err = scraper.GetPropertyDetails(propCtx, prop.URL)
				cancel()
			}
			if err != nil {
				// Manejo de error final
				log.Printf("❌ Error final al procesar propiedad: %v\n", err)
				fallidas++
				continue
			}
		}

		// Actualizar propiedad con los detalles
		prop.TipoPropiedad = ptr(details.TipoPropiedad)
		prop.Ubicacion = ptr(details.Ubicacion)
		prop.Dormitorios = ptr(details.Dormitorios)
		prop.Banios = ptr(details.Banios)
		prop.Antiguedad = ptr(details.Antiguedad)
		prop.SuperficieCubierta = ptr(details.SuperficieCubierta)
		prop.SuperficieTotal = ptr(details.SuperficieTotal)
		prop.SuperficieTerreno = ptr(details.SuperficieTerreno)
		prop.Frente = ptr(details.Frente)
		prop.Fondo = ptr(details.Fondo)
		prop.Ambientes = ptr(details.Ambientes)
		prop.Plantas = ptr(details.Plantas)
		prop.Cocheras = ptr(details.Cocheras)
		prop.Situacion = ptr(details.Situacion)
		prop.Expensas = ptr(details.Expensas)
		prop.Descripcion = ptr(details.Descripcion)
		prop.Status = "completed"
		actualizadas++

		prop.FechaScraping = time.Now()
		if err := database.UpdatePropiedadDetalles(&prop); err != nil {
			log.Printf("Error actualizando propiedad: %v\n", err)
			fallidas++
			continue
		}

		// Delay variable entre requests
		time.Sleep(time.Duration(2+rand.Intn(3)) * time.Second)
	}

	fmt.Printf("\nResumen:\n"+
		"- Propiedades actualizadas: %d\n"+
		"- Propiedades fallidas: %d\n"+
		"- Propiedades no disponibles: %d\n",
		actualizadas, fallidas, noDisponibles)

	return nil
}
