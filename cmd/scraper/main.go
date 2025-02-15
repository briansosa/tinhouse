package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/findhouse/internal/models"
	"github.com/findhouse/internal/scraper/tokko"
)

type Inmobiliaria struct {
	Nombre  string
	URL     string
	Sistema string
	Zona    string
}

func main() {
	// Agregar flag para modo test
	testMode := flag.Bool("test", false, "Ejecutar en modo test (solo una inmobiliaria)")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Cargar inmobiliarias desde CSV
	inmobiliarias, err := loadInmobiliarias("inmobiliarias_lanus.csv")
	if err != nil {
		log.Fatal("Error cargando inmobiliarias:", err)
	}

	// Filtrar solo las que usan Tokko
	var tokkoInmobiliarias []Inmobiliaria
	for _, inmo := range inmobiliarias {
		if strings.EqualFold(inmo.Sistema, "tokko") {
			tokkoInmobiliarias = append(tokkoInmobiliarias, inmo)
		}
	}

	if *testMode {
		// En modo test, usar solo la primera inmobiliaria
		if len(tokkoInmobiliarias) > 0 {
			tokkoInmobiliarias = tokkoInmobiliarias[:1]
			fmt.Println("Modo test: usando solo", tokkoInmobiliarias[0].Nombre)
		}
	}

	fmt.Printf("Encontradas %d inmobiliarias con Tokko\n", len(tokkoInmobiliarias))

	// Filtro común para todas
	filter := models.PropertyFilter{
		Operation:   "venta",
		Type:        "casa",
		Zone:        "G.B.A. Zona Sur",
		Location:    "Lanús",
		MinPriceUSD: 0,
		MaxPriceUSD: 80000,
	}

	// Leer propiedades existentes
	existingProperties := make(map[string]bool)
	if err := loadExistingProperties("properties.csv", existingProperties); err != nil {
		log.Printf("Advertencia al cargar propiedades existentes: %v\n", err)
	}

	// Preparar archivo CSV
	file, writer := prepareCSVFile("properties.csv")
	defer file.Close()
	defer writer.Flush()

	// Scrapear cada inmobiliaria
	for _, inmo := range tokkoInmobiliarias {
		fmt.Printf("\nScrapeando %s (%s)...\n", inmo.Nombre, inmo.URL)

		scraper := tokko.New(inmo.URL)
		properties, err := scraper.SearchProperties(ctx, filter)
		if err != nil {
			log.Printf("Error scrapeando %s: %v\n", inmo.Nombre, err)
			continue
		}

		// Guardar solo propiedades nuevas
		newCount := 0
		for _, prop := range properties {
			if !existingProperties[prop.Code] {
				if err := writeProperty(writer, prop, inmo.Nombre); err != nil {
					log.Printf("Error guardando propiedad %s: %v\n", prop.Code, err)
					continue
				}
				existingProperties[prop.Code] = true
				newCount++
			}
		}

		fmt.Printf("Encontradas %d propiedades en %s (%d nuevas)\n",
			len(properties), inmo.Nombre, newCount)
	}
}

func loadInmobiliarias(filename string) ([]Inmobiliaria, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	// Saltar headers
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error leyendo headers: %v", err)
	}

	// Encontrar índices de las columnas que nos interesan
	nombreIdx := findColumnIndex(headers, "Nombre")
	urlIdx := findColumnIndex(headers, "Sitio Web")
	sistemaIdx := findColumnIndex(headers, "Sistema")

	if nombreIdx == -1 || urlIdx == -1 || sistemaIdx == -1 {
		return nil, fmt.Errorf("no se encontraron todas las columnas necesarias")
	}

	var inmobiliarias []Inmobiliaria
	lineNum := 1

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Advertencia: error en línea %d: %v\n", lineNum, err)
			continue
		}

		if len(record) > sistemaIdx {
			nombre := strings.TrimSpace(record[nombreIdx])
			url := strings.TrimSpace(record[urlIdx])
			sistema := strings.TrimSpace(record[sistemaIdx])

			// Verificar que tengamos los datos necesarios
			if nombre != "" && url != "" && strings.Contains(strings.ToLower(sistema), "tokko") {
				inmobiliarias = append(inmobiliarias, Inmobiliaria{
					Nombre:  nombre,
					URL:     url,
					Sistema: "tokko",
					Zona:    "Lanús", // Por defecto
				})
			}
		}

		lineNum++
	}

	if len(inmobiliarias) == 0 {
		return nil, fmt.Errorf("no se encontraron inmobiliarias con Tokko")
	}

	return inmobiliarias, nil
}

func findColumnIndex(headers []string, name string) int {
	for i, h := range headers {
		if strings.EqualFold(strings.TrimSpace(h), name) {
			return i
		}
	}
	return -1
}

func prepareCSVFile(filename string) (*os.File, *csv.Writer) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error abriendo archivo CSV:", err)
	}

	writer := csv.NewWriter(file)
	writer.Comma = ';'

	// Escribir headers si el archivo está vacío
	if fi, err := file.Stat(); err == nil && fi.Size() == 0 {
		headers := []string{
			"Inmobiliaria",
			"Código",
			"Fecha Scraping",
			"Título",
			"Precio",
			"Dirección",
			"URL",
			"URL Imagen",
		}
		if err := writer.Write(headers); err != nil {
			log.Fatal("Error escribiendo headers:", err)
		}
	}

	return file, writer
}

func loadExistingProperties(filename string, properties map[string]bool) error {
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return nil // El archivo no existe, está bien
	}
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	// Saltar headers
	if _, err := reader.Read(); err != nil {
		return err
	}

	// Leer propiedades existentes
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) >= 2 {
			properties[record[1]] = true // El código está en la segunda columna
		}
	}

	return nil
}

func writeProperty(writer *csv.Writer, prop models.Property, inmobiliaria string) error {
	price := strings.Split(prop.PriceText, "\n")[0]
	price = strings.TrimSpace(price)
	title := strings.Join(strings.Fields(prop.Title), " ")

	row := []string{
		inmobiliaria,
		strings.TrimSpace(prop.Code),
		time.Now().Format("2006-01-02 15:04:05"),
		title,
		price,
		strings.TrimSpace(prop.Address),
		prop.URL,
		prop.ImageURL,
	}
	return writer.Write(row)
}
