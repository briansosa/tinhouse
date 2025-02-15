package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/findhouse/internal/db"
	"github.com/findhouse/internal/models"
	"github.com/findhouse/internal/scraper/tokko"
)

func main() {
	dbPath := flag.String("db", "findhouse.db", "Ruta a la base de datos SQLite")
	testMode := flag.Bool("test", false, "Ejecutar en modo test (solo una inmobiliaria)")
	flag.Parse()

	database, err := db.New(*dbPath)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer database.Close()

	// Por ahora usamos filtros fijos
	filter := models.PropertyFilter{
		Operation:   "venta",
		Type:        "casa",
		Zone:        "G.B.A. Zona Sur",
		Location:    "Lanús",
		MinPriceUSD: 0,
		MaxPriceUSD: 90000,
	}

	if err := searchProperties(database, filter, *testMode); err != nil {
		log.Fatal(err)
	}
}

func searchProperties(database *db.DB, filter models.PropertyFilter, testMode bool) error {
	// Obtener/crear la búsqueda
	busqueda, err := database.GetOrCreateBusqueda(filter)
	if err != nil {
		return fmt.Errorf("error preparando búsqueda: %v", err)
	}

	fmt.Printf("Búsqueda ID: %d (creada: %s)\n",
		busqueda.ID,
		busqueda.CreatedAt.Format("2006-01-02 15:04:05"))

	// Obtener todas las inmobiliarias con sistema
	inmobiliarias, err := database.GetInmobiliariasSistema()
	if err != nil {
		return fmt.Errorf("error obteniendo inmobiliarias: %v", err)
	}

	// Filtrar solo las que usan Tokko
	var tokkoInmobiliarias []db.Inmobiliaria
	for _, inmo := range inmobiliarias {
		if strings.Contains(strings.ToLower(inmo.Sistema), "tokko") {
			tokkoInmobiliarias = append(tokkoInmobiliarias, inmo)
		}
	}

	if testMode && len(tokkoInmobiliarias) > 0 {
		tokkoInmobiliarias = tokkoInmobiliarias[:1]
		fmt.Println("Modo test: usando solo", tokkoInmobiliarias[0].Nombre)
	}

	fmt.Printf("Encontradas %d inmobiliarias con sistema (%d usan Tokko)\n",
		len(inmobiliarias), len(tokkoInmobiliarias))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	for _, inmo := range tokkoInmobiliarias {
		fmt.Printf("\nScrapeando %s (%s)...\n", inmo.Nombre, inmo.URL)

		scraper := tokko.New(inmo.URL)
		properties, err := scraper.SearchProperties(ctx, filter)
		if err != nil {
			log.Printf("Error scrapeando %s: %v\n", inmo.Nombre, err)
			continue
		}

		var nuevas, actualizadas int
		for _, prop := range properties {
			// Convertir de models.Property a db.Propiedad
			propiedad := &db.Propiedad{
				InmobiliariaID: inmo.ID,
				Codigo:         prop.Code,
				Titulo:         prop.Title,
				Precio:         prop.PriceText,
				Direccion:      prop.Address,
				URL:            prop.URL,
				ImagenURL:      prop.ImageURL,
				FechaScraping:  time.Now(),
			}

			err := database.CreatePropiedadAndLink(propiedad, busqueda.ID)
			if err != nil {
				log.Printf("Error guardando propiedad %s: %v\n", propiedad.Codigo, err)
				continue
			}

			if propiedad.CreatedAt == propiedad.UpdatedAt {
				nuevas++
			} else {
				actualizadas++
			}
		}

		fmt.Printf("Encontradas %d propiedades en %s (%d nuevas, %d actualizadas)\n",
			len(properties), inmo.Nombre, nuevas, actualizadas)
	}

	return nil
}
