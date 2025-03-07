package analyzer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/findhouse/internal/db"
	"github.com/findhouse/internal/models"
	"github.com/findhouse/internal/scraper"
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

// SearchProperties busca propiedades en las inmobiliarias y las guarda en la DB
func SearchProperties(database *db.DB, testMode bool, inmobiliariaFilter string) error {
	ctx := context.Background()

	// Obtener inmobiliarias con sistema identificado
	inmobiliarias, err := database.GetInmobiliariasSistema()
	if err != nil {
		return fmt.Errorf("error obteniendo inmobiliarias: %v", err)
	}

	fmt.Printf("Encontradas %d inmobiliarias con sistema identificado\n", len(inmobiliarias))

	// Filtrar por nombre de inmobiliaria si se especificó
	var inmobiliariasFiltradas []db.Inmobiliaria
	if inmobiliariaFilter != "" {
		for _, inmo := range inmobiliarias {
			if strings.Contains(strings.ToLower(inmo.Nombre), strings.ToLower(inmobiliariaFilter)) {
				inmobiliariasFiltradas = append(inmobiliariasFiltradas, inmo)
			}
		}
		if len(inmobiliariasFiltradas) == 0 {
			return fmt.Errorf("no se encontró ninguna inmobiliaria con el nombre '%s'", inmobiliariaFilter)
		}
		inmobiliarias = inmobiliariasFiltradas
		fmt.Printf("Filtrando por inmobiliaria: %s (encontradas %d)\n", inmobiliariaFilter, len(inmobiliarias))
	}

	var totalPropiedades, nuevasPropiedades, propiedadesExistentes int

	var indexTest int
	for _, inmo := range inmobiliarias {
		fmt.Printf("\nScrapeando %s (%s)...\n", inmo.Nombre, inmo.URL)

		// Crear un scraper basado en el sistema de la inmobiliaria
		propertyScraper := scraper.NewScraper(inmo.Sistema, inmo.URL)
		if propertyScraper == nil {
			log.Printf("Sistema no soportado: %s\n", inmo.Sistema)
			continue
		}

		// Crear contexto con timeout para evitar bloqueos
		propCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)

		// Usar el scraper para buscar propiedades
		properties, err := propertyScraper.SearchProperties(propCtx)
		cancel()

		if err != nil {
			log.Printf("Error scrapeando %s: %v\n", inmo.Nombre, err)
			continue
		}

		fmt.Printf("Encontradas %d propiedades en %s\n", len(properties), inmo.Nombre)
		totalPropiedades += len(properties)

		// Procesar cada propiedad
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
				nuevasPropiedades++
			} else {
				propiedadesExistentes++
			}
		}

		indexTest++
		if testMode && indexTest == 5 {
			break
		}

		// Delay entre inmobiliarias para no sobrecargar
		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\nResumen:\n"+
		"- Total propiedades encontradas: %d\n"+
		"- Propiedades nuevas: %d\n"+
		"- Propiedades existentes: %d\n",
		totalPropiedades, nuevasPropiedades, propiedadesExistentes)

	return nil
}

// Función auxiliar para convertir a puntero
func ptr[T any](v T) *T {
	return &v
}

// normalizarTipoPropiedad normaliza el tipo de propiedad para que coincida con los códigos en la base de datos
func normalizarTipoPropiedad(tipo string) string {
	// Convertir a minúsculas y eliminar espacios
	tipo = strings.TrimSpace(tipo)

	// Mapeo de tipos comunes
	switch strings.ToLower(tipo) {
	case "casa", "chalet", "casa chalet", "casa quinta", "quinta":
		return "Casa"
	case "departamento", "depto", "dpto", "dpto.", "depto.", "departamento con dependencia":
		return "Departamento"
	case "ph", "p.h.", "p.h", "propiedad horizontal":
		return "PH"
	case "local", "local comercial", "fondo de comercio":
		return "Local"
	case "oficina", "consultorio", "estudio":
		return "Oficina"
	case "terreno", "lote", "lote de terreno", "fracción", "fraccion":
		return "Terreno"
	case "galpón", "galpon", "depósito", "deposito", "nave industrial":
		return "Galpón"
	default:
		// Si no hay coincidencia, devolvemos el tipo original
		return tipo
	}
}

// UpdateProperties actualiza los detalles de las propiedades en la base de datos
func UpdateProperties(database *db.DB, testMode bool, inmobiliariaFilter string) error {
	ctx := context.Background()

	// Obtener propiedades sin detalles
	propiedades, err := database.GetPropiedadesSinDetalles()
	if err != nil {
		return fmt.Errorf("error obteniendo propiedades sin detalles: %v", err)
	}

	fmt.Printf("Encontradas %d propiedades sin detalles\n", len(propiedades))

	// Filtrar por inmobiliaria si se especificó
	if inmobiliariaFilter != "" {
		var propiedadesFiltradas []db.Propiedad
		for _, prop := range propiedades {
			// Obtener información de la inmobiliaria
			inmobiliaria, err := database.GetInmobiliariaByID(prop.InmobiliariaID)
			if err != nil {
				log.Printf("Error obteniendo inmobiliaria %d: %v\n", prop.InmobiliariaID, err)
				continue
			}

			if strings.Contains(strings.ToLower(inmobiliaria.Nombre), strings.ToLower(inmobiliariaFilter)) {
				propiedadesFiltradas = append(propiedadesFiltradas, prop)
			}
		}

		if len(propiedadesFiltradas) == 0 {
			return fmt.Errorf("no se encontraron propiedades para la inmobiliaria '%s'", inmobiliariaFilter)
		}

		propiedades = propiedadesFiltradas
		fmt.Printf("Filtrando por inmobiliaria: %s (encontradas %d propiedades)\n", inmobiliariaFilter, len(propiedades))
	}

	// Definir el número de trabajadores concurrentes
	// Ajusta este valor según la capacidad de tu sistema
	numWorkers := 15
	if testMode {
		numWorkers = 2 // Menos trabajadores en modo de prueba
	}

	// Limitar el número de propiedades en modo de prueba
	if testMode && len(propiedades) > 5 {
		propiedades = propiedades[:5]
	}

	// Crear un canal para las propiedades a procesar
	propiedadesChan := make(chan db.Propiedad, len(propiedades))

	// Crear un canal para los resultados de extracción (no de escritura)
	extraccionesChan := make(chan struct {
		propiedad    db.Propiedad
		details      *models.PropertyDetails
		exito        bool
		noDisponible bool
	}, len(propiedades))

	// Crear un canal para los resultados finales
	resultadosChan := make(chan struct {
		exito        bool
		noDisponible bool
		propCodigo   string
	}, len(propiedades))

	// Crear un WaitGroup para esperar a que todos los trabajadores terminen
	var wg sync.WaitGroup

	// Iniciar los trabajadores de extracción
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			extraerPropiedades(ctx, database, propiedadesChan, extraccionesChan, workerID)
		}(i)
	}

	// Iniciar el trabajador de escritura (solo uno para evitar bloqueos)
	wg.Add(1)
	go func() {
		defer wg.Done()
		procesarEscrituras(database, extraccionesChan, resultadosChan)
	}()

	// Enviar las propiedades al canal
	for _, prop := range propiedades {
		propiedadesChan <- prop
	}
	close(propiedadesChan) // Cerrar el canal cuando todas las propiedades han sido enviadas

	// Crear un canal para recopilar los resultados
	doneChan := make(chan struct{})

	// Contadores para estadísticas
	var actualizadas, fallidas, noDisponibles int

	// Goroutine para recopilar resultados
	go func() {
		for resultado := range resultadosChan {
			if resultado.noDisponible {
				noDisponibles++
				fmt.Printf("⚠️ Propiedad %s no disponible\n", resultado.propCodigo)
			} else if resultado.exito {
				actualizadas++
				fmt.Printf("✓ Propiedad %s actualizada exitosamente\n", resultado.propCodigo)
			} else {
				fallidas++
				fmt.Printf("❌ Error actualizando propiedad %s\n", resultado.propCodigo)
			}
		}
		doneChan <- struct{}{}
	}()

	// Esperar a que todos los trabajadores terminen
	wg.Wait()
	close(extraccionesChan) // Cerrar el canal de extracciones
	close(resultadosChan)   // Cerrar el canal de resultados
	<-doneChan              // Esperar a que se procesen todos los resultados

	fmt.Printf("\nResumen:\n"+
		"- Propiedades actualizadas: %d\n"+
		"- Propiedades fallidas: %d\n"+
		"- Propiedades no disponibles: %d\n",
		actualizadas, fallidas, noDisponibles)

	return nil
}

// extraerPropiedades extrae los detalles de las propiedades sin escribir en la base de datos
func extraerPropiedades(ctx context.Context, database *db.DB, propiedadesChan <-chan db.Propiedad, extraccionesChan chan<- struct {
	propiedad    db.Propiedad
	details      *models.PropertyDetails
	exito        bool
	noDisponible bool
}, workerID int) {
	for prop := range propiedadesChan {
		fmt.Printf("\n[Worker %d] Extrayendo propiedad %s\n", workerID, prop.Codigo)
		fmt.Printf("   URL: %s\n", prop.URL)
		fmt.Printf("   Inmobiliaria ID: %d\n", prop.InmobiliariaID)

		resultado := struct {
			propiedad    db.Propiedad
			details      *models.PropertyDetails
			exito        bool
			noDisponible bool
		}{
			propiedad:    prop,
			details:      nil,
			exito:        false,
			noDisponible: false,
		}

		// Obtener información de la inmobiliaria
		inmobiliaria, err := database.GetInmobiliariaByID(prop.InmobiliariaID)
		if err != nil {
			log.Printf("[Worker %d] ❌ Error obteniendo inmobiliaria %d: %v\n", workerID, prop.InmobiliariaID, err)
			extraccionesChan <- resultado
			continue
		}

		fmt.Printf("   Inmobiliaria: %s (Sistema: %s)\n", inmobiliaria.Nombre, inmobiliaria.Sistema)

		// Obtener el scraper adecuado para el sistema de la inmobiliaria
		propertyScraper := scraper.NewScraper(inmobiliaria.Sistema, inmobiliaria.URL)
		if propertyScraper == nil {
			log.Printf("[Worker %d] Sistema no soportado: %s\n", workerID, inmobiliaria.Sistema)
			extraccionesChan <- resultado
			continue
		}

		// Obtener detalles usando el contexto con timeout
		propCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		log.Printf("[Worker %d] 🔍 Intentando obtener detalles de: %s\n", workerID, prop.URL)
		details, err := propertyScraper.GetPropertyDetails(propCtx, prop.URL)
		cancel() // Liberamos recursos

		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				log.Printf("[Worker %d] ⌛ Timeout al procesar propiedad: %s, reintentando...\n", workerID, prop.URL)
				// Reintento con más tiempo
				propCtx, cancel = context.WithTimeout(ctx, 5*time.Minute)
				details, err = propertyScraper.GetPropertyDetails(propCtx, prop.URL)
				cancel()
			}
			if err != nil {
				// Manejo de error final
				log.Printf("[Worker %d] ❌ Error final al procesar propiedad: %v\n", workerID, err)
				extraccionesChan <- resultado
				continue
			}
		}

		// Verificar si la propiedad ya no está disponible
		if details.Descripcion == "Propiedad no disponible" {
			log.Printf("[Worker %d] ⚠️ Propiedad no disponible: %s\n", workerID, prop.URL)
			resultado.noDisponible = true
			extraccionesChan <- resultado
			continue
		}

		log.Printf("[Worker %d] ✓ Extracción completada: %+v\n", workerID, details)

		// Si llegamos aquí, la extracción fue exitosa
		resultado.details = details
		resultado.exito = true
		extraccionesChan <- resultado

		// Delay fijo entre extracciones para evitar sobrecargar el sitio web
		time.Sleep(1 * time.Second)
	}
}

// procesarEscrituras procesa las escrituras en la base de datos de forma secuencial
func procesarEscrituras(database *db.DB, extraccionesChan <-chan struct {
	propiedad    db.Propiedad
	details      *models.PropertyDetails
	exito        bool
	noDisponible bool
}, resultadosChan chan<- struct {
	exito        bool
	noDisponible bool
	propCodigo   string
}) {
	for extraccion := range extraccionesChan {
		resultado := struct {
			exito        bool
			noDisponible bool
			propCodigo   string
		}{
			exito:        false,
			noDisponible: extraccion.noDisponible,
			propCodigo:   extraccion.propiedad.Codigo,
		}

		// Si la propiedad no está disponible o hubo un error en la extracción, solo enviamos el resultado
		if extraccion.noDisponible || !extraccion.exito {
			resultadosChan <- resultado
			continue
		}

		// Procesar los datos extraídos y guardarlos en la base de datos
		prop := extraccion.propiedad
		details := extraccion.details

		// Normalizar el tipo de propiedad
		tipoNormalizado := normalizarTipoPropiedad(details.TipoPropiedad)

		// Obtener el ID del tipo de propiedad
		var tipoID int64
		err := database.QueryRow(`
			SELECT id FROM property_types WHERE name = ?
		`, tipoNormalizado).Scan(&tipoID)
		if err != nil {
			// Si no existe, intentar insertar
			if err == sql.ErrNoRows {
				// Generar un código a partir del nombre
				code := strings.ToLower(strings.ReplaceAll(tipoNormalizado, " ", "_"))
				result, err := database.Exec(`
					INSERT INTO property_types (code, name) VALUES (?, ?)
				`, code, tipoNormalizado)
				if err != nil {
					log.Printf("[Escritor] Error al insertar tipo de propiedad '%s': %v\n", tipoNormalizado, err)
				} else {
					tipoID, _ = result.LastInsertId()
				}
			} else {
				log.Printf("[Escritor] Error al obtener ID para tipo de propiedad '%s': %v\n", tipoNormalizado, err)
			}
		}

		// Actualizar los campos de la propiedad con los detalles obtenidos
		prop.TipoPropiedad = &tipoID
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

		// Asignar coordenadas si están disponibles
		if details.Latitud != 0 && details.Longitud != 0 {
			prop.Latitud = &details.Latitud
			prop.Longitud = &details.Longitud
		}

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

		// Actualizar la propiedad en la base de datos
		log.Printf("[Escritor] Actualizando propiedad %s en la base de datos...\n", prop.Codigo)
		err = database.UpdatePropiedadDetalles(&prop)
		if err != nil {
			log.Printf("[Escritor] ❌ Error actualizando propiedad: %v\n", err)
			resultadosChan <- resultado
		} else {
			resultado.exito = true
			resultadosChan <- resultado
		}
	}
}
