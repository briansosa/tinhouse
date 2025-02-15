package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/findhouse/internal/analyzer"
	"github.com/findhouse/internal/db"
)

const (
	ModeFindInmobiliarias = "find-inmobiliarias"
	ModeAnalyzeSystems    = "analyze-systems"
	ModeNewInmobiliarias  = "new-inmobiliarias"
)

type Config struct {
	DBPath string
	Mode   string
}

func parseFlags() Config {
	config := Config{}
	flag.StringVar(&config.DBPath, "db", "findhouse.db", "Ruta a la base de datos SQLite")
	flag.StringVar(&config.Mode, "mode", ModeNewInmobiliarias, "Modo de ejecución (find-inmobiliarias, analyze-systems, new-inmobiliarias)")
	flag.Parse()
	return config
}

func validateMode(mode string) error {
	validModes := map[string]bool{
		ModeFindInmobiliarias: true,
		ModeAnalyzeSystems:    true,
		ModeNewInmobiliarias:  true,
	}

	if !validModes[mode] {
		return fmt.Errorf("modo inválido: %s", mode)
	}
	return nil
}

func printUsage() {
	fmt.Println("Modos válidos:")
	fmt.Printf("  - %s: solo busca nuevas inmobiliarias\n", ModeFindInmobiliarias)
	fmt.Printf("  - %s: solo analiza sistemas\n", ModeAnalyzeSystems)
	fmt.Printf("  - %s: ejecuta ambos procesos\n", ModeNewInmobiliarias)
}

func findInmobiliarias(database *db.DB) error {
	log.Println("Iniciando búsqueda de inmobiliarias...")
	return analyzer.SearchAndSaveInmobiliarias(database)
}

func analyzeSystems(database *db.DB) error {
	log.Println("Iniciando análisis de sistemas...")
	return analyzer.AnalyzeSitesDB(database)
}

func runProcess(database *db.DB, mode string) error {
	switch mode {
	case ModeFindInmobiliarias:
		return findInmobiliarias(database)

	case ModeAnalyzeSystems:
		return analyzeSystems(database)

	case ModeNewInmobiliarias:
		if err := findInmobiliarias(database); err != nil {
			return fmt.Errorf("error en búsqueda de inmobiliarias: %w", err)
		}

		if err := analyzeSystems(database); err != nil {
			return fmt.Errorf("error en análisis de sistemas: %w", err)
		}
	}
	return nil
}

func main() {
	config := parseFlags()

	if err := validateMode(config.Mode); err != nil {
		fmt.Printf("Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}

	database, err := db.New(config.DBPath)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer database.Close()

	if err := runProcess(database, config.Mode); err != nil {
		log.Fatal(err)
	}

	log.Printf("Proceso '%s' completado exitosamente\n", config.Mode)
}
