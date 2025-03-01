package analyzer

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type SistemaInmobiliario struct {
	Nombre      string
	Marcadores  []string
	Ocurrencias int
}

var sistemasConocidos = []SistemaInmobiliario{
	{
		Nombre: "Tokko Broker",
		Marcadores: []string{
			"tokkobroker",
			"tkb.com.ar",
		},
	},
	{
		Nombre: "Properati",
		Marcadores: []string{
			"properati",
			"property-gallery",
		},
	},
	{
		Nombre: "Zonaprop",
		Marcadores: []string{
			"zonaprop.com.ar",
		},
	},
	{
		Nombre: "WordPress",
		Marcadores: []string{
			"wp-content",
			"wp-includes",
		},
	},
	{
		Nombre: "Buscador Prop",
		Marcadores: []string{
			"buscadorprop",
			"grupotodo.com.ar",
		},
	},
	{
		Nombre: "Argencasas",
		Marcadores: []string{
			"argencasas.com",
			"argencasas.com.ar",
			"inmobiliario.com.ar",
		},
	},
	{
		Nombre: "Ubiquo",
		Marcadores: []string{
			"ubiquo",
			"ubiquo.com.ar",
		},
	},
	{
		Nombre: "Adinco",
		Marcadores: []string{
			"adinco",
			"crm.adinco.net",
		},
	},
	{
		Nombre: "Me Mudo Ya",
		Marcadores: []string{
			"memudoya",
			"memudoya.com",
			"mapaprop",
			"mapaprop.com",
		},
	},
	{
		Nombre: "Amaira",
		Marcadores: []string{
			"amaira",
			"amaira.com.ar",
			"xintel",
			"xintel.com.ar",
		},
	},
	{
		Nombre: "Desarrollo propio",
		Marcadores: []string{
			"sysmika",
			"fenix",
			"nibiru",
		},
	},
}

func AnalyzeSites() {
	// Inicializar Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	file, err := os.Open("inmobiliarias_lanus.csv")
	if err != nil {
		log.Fatal("Error al abrir CSV:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	inmobiliarias, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error al leer CSV:", err)
	}

	// Mapa para contar sistemas
	resultados := make(map[string][]string)
	noIdentificados := make([]string, 0)

	// Verificar si existe la columna Sistema
	headers := inmobiliarias[0]
	sistemaIndex := -1
	for i, header := range headers {
		if header == "Sistema" {
			sistemaIndex = i
			break
		}
	}

	if sistemaIndex == -1 {
		log.Fatal("No se encontró la columna Sistema en el CSV")
	}

	// Procesar solo las inmobiliarias sin sistema o con "No identificado"
	for _, record := range inmobiliarias[1:] {
		if len(record) <= sistemaIndex || record[sistemaIndex] == "" || record[sistemaIndex] == "No identificado" {
			nombre := record[0]
			sitioWeb := record[4]

			if sitioWeb == "" {
				continue
			}

			fmt.Printf("Analizando %s (%s)...\n", nombre, sitioWeb)

			var htmlContent string
			err := chromedp.Run(ctx,
				chromedp.Navigate(sitioWeb),
				chromedp.Sleep(2*time.Second),
				chromedp.OuterHTML("html", &htmlContent),
			)

			if err != nil {
				log.Printf("Error al acceder a %s: %v", sitioWeb, err)
				continue
			}

			sistemaEncontrado := false
			for _, sistema := range sistemasConocidos {
				for _, marcador := range sistema.Marcadores {
					if strings.Contains(strings.ToLower(htmlContent), strings.ToLower(marcador)) {
						resultados[sistema.Nombre] = append(resultados[sistema.Nombre], fmt.Sprintf("%s (%s)", nombre, sitioWeb))
						sistemaEncontrado = true
						break
					}
				}
				if sistemaEncontrado {
					break
				}
			}

			if !sistemaEncontrado {
				noIdentificados = append(noIdentificados, fmt.Sprintf("%s (%s)", nombre, sitioWeb))
			}
		}
	}

	// Guardar resultados
	if err := guardarResultados(inmobiliarias, resultados); err != nil {
		log.Fatal("Error al guardar resultados:", err)
	}

	// Mostrar resumen de nuevos hallazgos
	fmt.Println("\n=== Nuevos Sistemas Identificados ===")
	for sistema, sitios := range resultados {
		fmt.Printf("\n%s (%d sitios):\n", sistema, len(sitios))
		for _, sitio := range sitios {
			fmt.Printf("  - %s\n", sitio)
		}
	}

	if len(noIdentificados) > 0 {
		fmt.Printf("\nSiguen sin identificar (%d sitios):\n", len(noIdentificados))
		for _, sitio := range noIdentificados {
			fmt.Printf("  - %s\n", sitio)
		}
	}
}

// AnalyzeSystem analiza una URL y detecta qué sistema usa
func AnalyzeSystem(url string) (string, error) {
	if url == "" {
		return "No identificado", nil
	}

	// Limpiar URL
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}

	// Configurar Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Agregar timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.OuterHTML("html", &htmlContent),
	)

	if err != nil {
		return "", fmt.Errorf("error al acceder a %s: %v", url, err)
	}

	// Buscar sistema usando los marcadores conocidos
	for _, sistema := range sistemasConocidos {
		for _, marcador := range sistema.Marcadores {
			if strings.Contains(strings.ToLower(htmlContent), strings.ToLower(marcador)) {
				return sistema.Nombre, nil
			}
		}
	}

	return "No identificado", nil
}

func guardarResultados(inmobiliarias [][]string, sistemas map[string][]string) error {
	sistemaPorInmobiliaria := make(map[string]string)

	for sistema, sitios := range sistemas {
		for _, sitio := range sitios {
			nombre := strings.Split(sitio, " (")[0]
			sistemaPorInmobiliaria[nombre] = sistema
		}
	}

	file, err := os.Create("inmobiliarias_lanus.csv")
	if err != nil {
		return fmt.Errorf("error al crear archivo: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir los encabezados originales
	if err := writer.Write(inmobiliarias[0]); err != nil {
		return fmt.Errorf("error al escribir encabezados: %v", err)
	}

	// Actualizar la columna Sistema existente
	for i, row := range inmobiliarias {
		if i == 0 { // Saltar encabezados
			continue
		}

		nombre := row[0]
		if sistema := sistemaPorInmobiliaria[nombre]; sistema != "" {
			// Actualizar el sistema en la columna existente (índice 5)
			row[5] = sistema
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error al escribir fila: %v", err)
		}
	}

	return nil
}
