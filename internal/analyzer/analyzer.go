package analyzer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/findhouse/internal/db"
	"github.com/findhouse/internal/scraper"
	"github.com/findhouse/internal/scraper/tokko"
	"golang.org/x/exp/rand"
)

// SearchAndSaveInmobiliarias busca inmobiliarias en Google Maps y guarda solo las nuevas en la DB
func SearchAndSaveInmobiliarias(database *db.DB, zone string) error {
	ctx := context.Background()

	// Usar el scraper existente para buscar inmobiliarias
	results, err := scraper.SearchInmobiliarias(ctx, zone)
	if err != nil {
		return fmt.Errorf("error buscando inmobiliarias: %v", err)
	}

	fmt.Printf("Encontradas %d inmobiliarias en Google Maps para la zona: %s\n", len(results), zone)

	var nuevas, existentes, dudosas int
	var nuevasLista []string
	var dudosasLista []string

	// Guardar cada inmobiliaria en la DB solo si no existe
	for _, result := range results {
		// Verificar si ya existe
		exists, err := database.ExistsInmobiliaria(result.Nombre, result.Direccion)
		if err != nil {
			continue
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
			Zona:      zone,
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

// AnalyzeSystem analiza las inmobiliarias y guarda/actualiza en la base de datos
func AnalyzeSystem(database *db.DB) error {
	// Obtener inmobiliarias sin sistema identificado
	inmobiliarias, err := database.GetInmobiliariasSinSistema()
	if err != nil {
		return fmt.Errorf("error obteniendo inmobiliarias: %v", err)
	}

	fmt.Printf("Analizando %d inmobiliarias...\n", len(inmobiliarias))

	for _, inmo := range inmobiliarias {
		system, err := scraper.AnalyzeSystem(inmo.URL)
		if err != nil {
			log.Printf("Error detectando sistema para %s: %v\n", inmo.Nombre, err)
			continue
		}

		// Actualizar inmobiliaria con el sistema detectado
		inmo.Sistema = system
		inmo.UpdatedAt = time.Now()

		if err := database.UpdateInmobiliariaSistema(&inmo); err != nil {
			log.Printf("Error actualizando sistema para %s: %v\n", inmo.Nombre, err)
			continue

		}

		fmt.Printf("✓ %s: %s\n", inmo.Nombre, system)
	}

	return nil
}

func SearchProperties(database *db.DB, testMode bool) error {
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
		// Ya no pasamos filtros al scraper
		properties, err := scraper.SearchProperties(ctx)
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
				Moneda:         prop.Currency,
				Direccion:      prop.Address,
				URL:            prop.URL,
				ImagenURL:      prop.ImageURL,
				Status:         "pending",
			}

			// Ya no vinculamos con búsqueda
			err := database.CreatePropiedad(propiedad)
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

		// Actualizar la propiedad con los detalles
		prop.TipoPropiedad = &details.TipoPropiedad
		prop.Ubicacion = &details.Ubicacion
		prop.Dormitorios = &details.Dormitorios
		prop.Banios = &details.Banios
		prop.Antiguedad = &details.Antiguedad
		prop.SuperficieCubierta = &details.SuperficieCubierta
		prop.SuperficieTotal = &details.SuperficieTotal
		prop.SuperficieTerreno = &details.SuperficieTerreno
		prop.Frente = &details.Frente
		prop.Fondo = &details.Fondo
		prop.Ambientes = &details.Ambientes
		prop.Plantas = &details.Plantas
		prop.Cocheras = &details.Cocheras
		prop.Situacion = &details.Situacion
		prop.Expensas = &details.Expensas
		prop.Descripcion = &details.Descripcion
		prop.Imagenes = &details.Images
		prop.Operacion = &details.Operacion
		prop.Condicion = &details.Condicion
		prop.Orientacion = &details.Orientacion
		prop.Disposicion = &details.Disposicion

		// Preparar las características para guardar
		prop.Features = make(map[string][]string)
		if len(details.Servicios) > 0 {
			prop.Features["servicio"] = details.Servicios
		}
		if len(details.TiposAmbientes) > 0 {
			prop.Features["ambiente"] = details.TiposAmbientes
		}
		if len(details.Adicionales) > 0 {
			prop.Features["adicional"] = details.Adicionales
		}

		prop.Status = "completed"
		actualizadas++

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
