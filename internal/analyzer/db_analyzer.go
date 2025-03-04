package analyzer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/findhouse/internal/db"
	"github.com/findhouse/internal/models"
	"github.com/findhouse/internal/scraper"
	"github.com/findhouse/internal/scraper/tokko"
	"golang.org/x/exp/rand"
)

// SearchAndSaveInmobiliarias busca inmobiliarias en Google Maps y guarda solo las nuevas en la DB
func SearchAndSaveInmobiliarias(database *db.DB) error {
	ctx := context.Background()

	// Usar el scraper existente para buscar inmobiliarias
	results, err := scraper.SearchInmobiliarias(ctx, "Lanús")
	if err != nil {
		return fmt.Errorf("error buscando inmobiliarias: %v", err)
	}

	fmt.Printf("Encontradas %d inmobiliarias en Google Maps\n", len(results))

	var nuevas, existentes, dudosas int
	var nuevasLista []string
	var dudosasLista []string

	// Guardar cada inmobiliaria en la DB solo si no existe
	for _, result := range results {
		// Verificar si ya existe
		log.Printf("Verificando existencia de result %v", result)
		exists, err := database.ExistsInmobiliaria(result.Nombre, result.Direccion)
		if err != nil {
			log.Printf("Error verificando existencia de %s: %v\n", result.Nombre, err)
			return err

		}

		if exists {
			existentes++
			continue
		}

		// Si la dirección está vacía, podría ser un caso dudoso
		if result.Direccion == "" {
			log.Printf("⚠️ Caso dudoso (sin dirección): %s\n", result.Nombre)
			dudosas++
			dudosasLista = append(dudosasLista, result.Nombre)
		}

		// Convertir rating de string a float64
		rating, err := strconv.ParseFloat(result.Rating, 64)
		if err != nil {
			rating = 0 // Valor por defecto si hay error en la conversión
		}

		inmo := &db.Inmobiliaria{
			Nombre:    result.Nombre,
			URL:       result.SitioWeb,
			Direccion: result.Direccion,
			Telefono:  result.Telefono,
			Rating:    rating,
			Zona:      "Lanús",
		}

		if err := database.CreateInmobiliaria(inmo); err != nil {
			log.Printf("Error guardando inmobiliaria %s: %v\n", inmo.Nombre, err)
			continue
		}
		nuevas++
		nuevasLista = append(nuevasLista, fmt.Sprintf("%s (%s)", inmo.Nombre, inmo.Direccion))
	}

	fmt.Printf("\nResumen:\n"+
		"- Nuevas inmobiliarias guardadas: %d\n"+
		"- Inmobiliarias existentes: %d\n"+
		"- Casos dudosos: %d\n", nuevas, existentes, dudosas)

	if len(nuevasLista) > 0 {
		fmt.Println("\nNuevas inmobiliarias:")
		for _, inmo := range nuevasLista {
			fmt.Printf("✅ %s\n", inmo)
		}
	}

	if len(dudosasLista) > 0 {
		fmt.Println("\nCasos dudosos:")
		for _, inmo := range dudosasLista {
			fmt.Printf("⚠️ %s\n", inmo)
		}
	}

	return nil
}

// AnalyzeSitesDB analiza las inmobiliarias y guarda/actualiza en la base de datos
func AnalyzeSitesDB(database *db.DB) error {
	// Obtener inmobiliarias sin sistema identificado
	inmobiliarias, err := database.GetInmobiliariasSinSistema()
	if err != nil {
		return fmt.Errorf("error obteniendo inmobiliarias: %v", err)
	}

	fmt.Printf("Analizando %d inmobiliarias...\n", len(inmobiliarias))

	for _, inmo := range inmobiliarias {
		sistema, err := DetectarSistema(inmo.URL)
		if err != nil {
			log.Printf("Error detectando sistema para %s: %v\n", inmo.Nombre, err)
			continue
		}

		// Actualizar inmobiliaria con el sistema detectado
		inmo.Sistema = sistema
		inmo.UpdatedAt = time.Now()

		if err := database.UpdateInmobiliariaSistema(&inmo); err != nil {
			log.Printf("Error actualizando sistema para %s: %v\n", inmo.Nombre, err)
			continue

		}

		fmt.Printf("✓ %s: %s\n", inmo.Nombre, sistema)
	}

	return nil
}

// DetectarSistema es la misma función que ya tienes pero exportada
func DetectarSistema(url string) (string, error) {
	// Implementar la lógica de detección que ya tienes en analyzer.go
	system, err := AnalyzeSystem(url)
	return system, err
}

func SearchProperties(database *db.DB, filter models.PropertyFilter, testMode bool) error {
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

// Función auxiliar para convertir a puntero
func ptr[T any](v T) *T {
	return &v
}

func UpdateProperties(database *db.DB, testMode bool) error {
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

		// Log para verificar los detalles obtenidos
		fmt.Printf("   Detalles obtenidos: %+v\n", details)

		// Actualizar propiedad con los detalles
		fmt.Printf("   Actualizando propiedad %s con detalles...\n", prop.Codigo)
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
		prop.Imagenes = ptr(details.Images)
		prop.Status = "completed"
		actualizadas++

		prop.FechaScraping = time.Now()
		if err := database.UpdatePropiedadDetalles(&prop); err != nil {
			log.Printf("Error actualizando propiedad: %v\n", err)
			fallidas++
			continue
		}

		// Log para confirmar la actualización
		fmt.Printf("✓ Propiedad %s actualizada exitosamente\n", prop.Codigo)

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
