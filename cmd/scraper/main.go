package main

import (
	"fmt"
	"log"
	"os"

	"github.com/findhouse/cmd/configuration"
	"github.com/findhouse/internal/analyzer"
	"github.com/findhouse/internal/db"
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

	if err := runProcess(database, flags); err != nil {
		log.Fatal(err)
	}

	log.Printf("Proceso '%s' completado exitosamente\n", flags.Mode)
}

func runProcess(database *db.DB, flags *configuration.Flags) error {
	switch flags.Mode {
	case configuration.ModeFindInmobiliarias:
		if flags.Zone == "" {
			return fmt.Errorf("zona no especificada")
		}

		return findInmobiliarias(database, flags.Zone)

	case configuration.ModeAnalyzeSystems:
		return analyzeSystems(database)

	case configuration.ModeNewInmobiliarias:
		if flags.Zone == "" {
			return fmt.Errorf("zona no especificada")
		}

		if err := findInmobiliarias(database, flags.Zone); err != nil {
			return fmt.Errorf("error en búsqueda de inmobiliarias: %w", err)
		}

		if err := analyzeSystems(database); err != nil {
			return fmt.Errorf("error en análisis de sistemas: %w", err)
		}

	case configuration.ModeSearchProperties:
		// Eliminamos los filtros fijos
		if err := searchProperties(database, flags.TestMode); err != nil {
			return fmt.Errorf("error en búsqueda de propiedades: %w", err)
		}

	case configuration.ModeUpdateProperties:
		if err := updateProperties(database, flags.TestMode); err != nil {
			return fmt.Errorf("error en actualización de propiedades: %w", err)
		}
	default:
		return fmt.Errorf("modo no válido: %s", flags.Mode)
	}

	return nil
}

func findInmobiliarias(database *db.DB, zone string) error {
	log.Println("Iniciando búsqueda de inmobiliarias...")
	return analyzer.SearchAndSaveInmobiliarias(database, zone)
}

func analyzeSystems(database *db.DB) error {
	log.Println("Iniciando análisis de sistemas...")
	return analyzer.AnalyzeSystem(database)
}

func searchProperties(database *db.DB, testMode bool) error {
	return analyzer.SearchProperties(database, testMode)
}

func updateProperties(database *db.DB, testMode bool) error {
	return analyzer.UpdateProperties(database, testMode)
}
