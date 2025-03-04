package main

import (
	"fmt"
	"log"
	"os"

	"github.com/findhouse/cmd/configuration"
	"github.com/findhouse/internal/analyzer"
	"github.com/findhouse/internal/db"
	"github.com/findhouse/internal/models"
)

func main() {
	flags, err := configuration.ParseFlags()
	if err != nil {
		log.Fatal("Error parsing flags:", err)
		os.Exit(1)
	}

	database, err := db.New(flags.DBPath)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := runProcess(database, string(flags.Mode), flags.TestMode); err != nil {
		log.Fatal(err)
	}

	log.Printf("Proceso '%s' completado exitosamente\n", flags.Mode)
}

func runProcess(database *db.DB, mode string, testMode bool) error {
	switch mode {
	case string(configuration.ModeFindInmobiliarias):
		return findInmobiliarias(database)

	case string(configuration.ModeAnalyzeSystems):
		return analyzeSystems(database)

	case string(configuration.ModeNewInmobiliarias):
		if err := findInmobiliarias(database); err != nil {
			return fmt.Errorf("error en búsqueda de inmobiliarias: %w", err)
		}

		if err := analyzeSystems(database); err != nil {
			return fmt.Errorf("error en análisis de sistemas: %w", err)
		}

	case string(configuration.ModeSearchProperties):
		// Por ahora usamos filtros fijos
		filter := models.PropertyFilter{
			Operation:   "venta",
			Type:        "casa",
			Zone:        "G.B.A. Zona Sur",
			Location:    "Lanús",
			MinPriceUSD: 0,
			MaxPriceUSD: 90000,
		}

		if err := searchProperties(database, filter, testMode); err != nil {
			return fmt.Errorf("error en búsqueda de propiedades: %w", err)
		}

	case string(configuration.ModeUpdateProperties):
		if err := updateProperties(database, testMode); err != nil {
			return fmt.Errorf("error en actualización de propiedades: %w", err)
		}
	default:
		return fmt.Errorf("modo no válido: %s", mode)
	}

	return nil
}

func findInmobiliarias(database *db.DB) error {
	log.Println("Iniciando búsqueda de inmobiliarias...")
	return analyzer.SearchAndSaveInmobiliarias(database)
}

func analyzeSystems(database *db.DB) error {
	log.Println("Iniciando análisis de sistemas...")
	return analyzer.AnalyzeSitesDB(database)
}

func searchProperties(database *db.DB, filter models.PropertyFilter, testMode bool) error {
	return analyzer.SearchProperties(database, filter, testMode)
}

func updateProperties(database *db.DB, testMode bool) error {
	return analyzer.UpdateProperties(database, testMode)
}
