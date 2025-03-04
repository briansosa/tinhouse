package configuration

import (
	"flag"
	"fmt"
	"os"
)

type ExecutionMode string

const (
	ModeFindInmobiliarias ExecutionMode = "find-inmobiliarias"
	ModeAnalyzeSystems    ExecutionMode = "analyze-systems"
	ModeNewInmobiliarias  ExecutionMode = "new-inmobiliarias"
	ModeSearchProperties  ExecutionMode = "search-properties"
	ModeUpdateProperties  ExecutionMode = "update-properties"
)

type Flags struct {
	Mode     ExecutionMode
	DBPath   string
	TestMode bool
	Zone     string // Zona para búsqueda de inmobiliarias
}

func ParseFlags() (*Flags, error) {
	flags := &Flags{}
	mode := string(flags.Mode)
	flag.StringVar(&mode, "mode", string(ModeNewInmobiliarias), "Modo de ejecución")
	flag.StringVar(&flags.DBPath, "db", "internal/db/findhouse.db", "Ruta a la base de datos SQLite")
	flag.BoolVar(&flags.TestMode, "test", false, "Ejecutar en modo de prueba")
	flag.StringVar(&flags.Zone, "zone", "", "Zona para búsqueda de inmobiliarias (ej: Lanús, Avellaneda, etc)")

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("error parsing flags: %w", err)
	}

	flags.Mode = ExecutionMode(mode)

	return flags, nil
}
