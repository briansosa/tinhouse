package tokko

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/findhouse/internal/models"
)

type Scraper struct {
	BaseURL string
}

func New(baseURL string) *Scraper {
	cleanURL := strings.TrimRight(baseURL, "/")
	return &Scraper{
		BaseURL: cleanURL,
	}
}

func (s *Scraper) SearchProperties(ctx context.Context) ([]models.Property, error) {
	baseURL := strings.TrimRight(s.BaseURL, "/")

	// Simplificamos la URL para obtener todas las propiedades sin filtros
	url := fmt.Sprintf("%s/Buscar", baseURL)

	fmt.Printf("Buscando en URL: %s\n", url)

	var properties []models.Property

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(ctx, opts...)
	defer allocCancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(5*time.Second),
		chromedp.Evaluate(`
			(() => {
				console.log('=== Debug Info ===');
				
				// Buscar el contenedor de propiedades
				const container = document.querySelector('#propiedades.resultados-list');
				console.log('Container found:', container ? 'yes' : 'no');
				
				// Buscar los items de propiedades (li elements dentro del contenedor)
				const properties = container ? Array.from(container.querySelectorAll('li')) : [];
				console.log('Properties found:', properties.length);
				
				return properties.map(prop => {
					// Extraer datos usando las clases que vemos en el DOM
					const priceElement = prop.querySelector('.prop-valor-nro');
					const codeElement = prop.querySelector('.codref');
					const titleElement = prop.querySelector('.prop-desc-tipo-ub');
					const addressElement = prop.querySelector('.prop-desc-dir');
					const imageElement = prop.querySelector('.dest-img');
					
					// Procesar el precio para separar moneda y valor
					let priceText = priceElement?.textContent?.trim() || '';
					
					// Limpiar el precio de cualquier código o texto adicional
					// Solo queremos mantener los dígitos, puntos y comas
					let currency = '';
					let priceValue = '';
					
					// Extraer moneda (USD o $)
					if (priceText.includes('USD')) {
						currency = 'USD';
						// Eliminar 'USD' y luego limpiar cualquier texto que no sea parte del precio
						priceValue = priceText.replace('USD', '').trim();
						// Usar una expresión regular para extraer solo el formato de precio (números, puntos, comas)
						const priceMatch = priceValue.match(/[\d.,]+/);
						priceValue = priceMatch ? priceMatch[0] : priceValue;
					} else if (priceText.includes('$')) {
						currency = 'ARS';
						// Eliminar '$' y luego limpiar cualquier texto que no sea parte del precio
						priceValue = priceText.replace('$', '').trim();
						// Usar una expresión regular para extraer solo el formato de precio (números, puntos, comas)
						const priceMatch = priceValue.match(/[\d.,]+/);
						priceValue = priceMatch ? priceMatch[0] : priceValue;
					} else {
						currency = 'Desconocida';
						// Intentar extraer solo el formato de precio (números, puntos, comas)
						const priceMatch = priceText.match(/[\d.,]+/);
						priceValue = priceMatch ? priceMatch[0] : priceText;
					}
					
					// Extraer el código correctamente, asegurándose de que no se mezcle con el precio
					let code = codeElement?.textContent?.trim() || '';
					
					return {
						title: titleElement?.textContent?.trim(),
						priceText: priceValue,
						currency: currency,
						address: addressElement?.textContent?.trim(),
						code: code,
						url: prop.querySelector('a')?.href,
						imageUrl: imageElement?.src || ''
					};
				});
			})()
		`, &properties),
	)

	if err != nil {
		return nil, fmt.Errorf("error extrayendo propiedades: %v", err)
	}

	return properties, nil
}

// GetPropertyDetails obtiene los detalles de una propiedad específica
func (s *Scraper) GetPropertyDetails(ctx context.Context, url string) (*models.PropertyDetails, error) {
	fmt.Printf("🔍 Intentando obtener detalles de: %s\n", url)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-notifications", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var details models.PropertyDetails

	err := chromedp.Run(taskCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("#ficha_detalle_cuerpo", chromedp.ByID),
		chromedp.Evaluate(`
			(() => {
				// Función auxiliar para extraer números
				function extractNumber(text) {
					if (!text) return 0;
					const match = text.match(/[\d,.]+/);
					return match ? parseFloat(match[0].replace(',', '.')) : 0;
				}

				// Función para extraer metros cuadrados
				function extractM2(text) {
					if (!text) return 0;
					const match = text.match(/(\d+(?:[,.]\d+)?)\s*m²/);
					return match ? parseFloat(match[1].replace(',', '.')) : 0;
				}

				// Función para buscar valor en la lista
				function findValue(selector, label) {
					const items = document.querySelectorAll(selector + ' li');
					for (const item of items) {
						if (item.textContent.toLowerCase().includes(label.toLowerCase())) {
							const text = item.textContent;
							const colonIndex = text.indexOf(':');
							return colonIndex !== -1 ? text.substring(colonIndex + 1).trim() : '';
						}
					}
					return '';
				}

				// Función para buscar texto en ficha_detalle_item
				function findDetailValue(label) {
					const items = document.querySelectorAll('#ficha_detalle_cuerpo .ficha_detalle_item');
					for (const item of items) {
						if (item.textContent.toLowerCase().includes(label.toLowerCase())) {
							return item.textContent.split(label)[1]?.trim() || '';
						}
					}
					return '';
				}

				// Extraer datos básicos
				const tipoPropiedad = document.querySelector('#ficha_detalle_cuerpo .ficha_detalle_item:first-child')?.textContent.split('Tipo de Propiedad')[1]?.trim() || '';
				const ubicacion = document.querySelector('#ficha_detalle_cuerpo .ficha_detalle_item:nth-child(2)')?.textContent.split('Ubicación')[1]?.trim() || '';
				
				// Extraer tipo de operación
				let operacion = '';
				
				// Método 1: Buscar en el título de la página
				const pageTitle = document.title || '';
				if (pageTitle.toLowerCase().includes('venta')) {
					operacion = 'Venta';
				} else if (pageTitle.toLowerCase().includes('alquiler temp')) {
					operacion = 'Alquiler Temporario';
				} else if (pageTitle.toLowerCase().includes('alquiler')) {
					operacion = 'Alquiler';
				}
				
				// Método 2: Buscar en la URL
				if (!operacion) {
					const currentUrl = window.location.href;
					if (currentUrl.toLowerCase().includes('venta')) {
						operacion = 'Venta';
					} else if (currentUrl.toLowerCase().includes('alquiler-temp')) {
						operacion = 'Alquiler Temporario';
					} else if (currentUrl.toLowerCase().includes('alquiler')) {
						operacion = 'Alquiler';
					}
				}
				
				// Método 3: Buscar en el contenido de la página
				if (!operacion) {
					const breadcrumbs = document.querySelector('.breadcrumb');
					if (breadcrumbs) {
						const breadcrumbText = breadcrumbs.textContent.toLowerCase();
						if (breadcrumbText.includes('venta')) {
							operacion = 'Venta';
						} else if (breadcrumbText.includes('alquiler temp')) {
							operacion = 'Alquiler Temporario';
						} else if (breadcrumbText.includes('alquiler')) {
							operacion = 'Alquiler';
						}
					}
				}
				
				// Extraer información básica
				const dormitorios = extractNumber(findValue('#lista_informacion_basica', 'Dormitorios'));
				const banios = extractNumber(findValue('#lista_informacion_basica', 'Baños'));
				const antiguedad = findValue('#lista_informacion_basica', 'Antigüedad');
				const ambientes = extractNumber(findValue('#lista_informacion_basica', 'Ambientes'));
				const plantas = extractNumber(findValue('#lista_informacion_basica', 'Plantas'));
				const cocheras = extractNumber(findValue('#lista_informacion_basica', 'Cocheras'));
				const situacion = findValue('#lista_informacion_basica', 'Situación');
				const expensas = extractNumber(findValue('#lista_informacion_basica', 'Expensas'));
				const condicion = findValue('#lista_informacion_basica', 'Condición') || findDetailValue('Condición');
				const orientacion = findValue('#lista_informacion_basica', 'Orientación') || findDetailValue('Orientación');
				const disposicion = findValue('#lista_informacion_basica', 'Disposición') || findDetailValue('Disposición');

				// Extraer superficies
				const superficieTerreno = extractM2(findValue('#lista_superficies', 'Terreno'));
				const superficieTotal = extractM2(findValue('#lista_superficies', 'Superficie Total'));
				const superficieCubierta = extractM2(findValue('#lista_superficies', 'Cubierta')) || 
					extractM2(findDetailValue('Total construido'));
				const frente = extractNumber(findValue('#lista_superficies', 'Frente'));
				const fondo = extractNumber(findValue('#lista_superficies', 'Fondo'));

				const descripcion = document.querySelector('#prop-desc')?.textContent.trim() || '';

				// Extraer imágenes
				const images = Array.from(document.querySelectorAll('#ficha_slider .slides li:not(.bx-clone) img'))
					.map(img => img.src);

				// Extraer coordenadas del mapa
				let latitud = 0;
				let longitud = 0;

				// Método 1: Buscar en el iframe de Google Maps
				const mapIframe = document.querySelector('iframe[src*="google.com/maps"]');
				if (mapIframe) {
					const src = mapIframe.src;
					const latMatch = src.match(/q=(-?\d+\.\d+),/);
					const lngMatch = src.match(/,(-?\d+\.\d+)/);
					
					if (latMatch && latMatch[1]) {
						latitud = parseFloat(latMatch[1]);
					}
					
					if (lngMatch && lngMatch[1]) {
						longitud = parseFloat(lngMatch[1]);
					}
				}

				// Método 2: Buscar en atributos data-* de elementos del mapa
				if (latitud === 0 && longitud === 0) {
					const mapElement = document.querySelector('[data-lat][data-lng]');
					if (mapElement) {
						const lat = mapElement.getAttribute('data-lat');
						const lng = mapElement.getAttribute('data-lng');
						
						if (lat) latitud = parseFloat(lat);
						if (lng) longitud = parseFloat(lng);
					}
				}

				// Método 3: Buscar en scripts de la página
				if (latitud === 0 && longitud === 0) {
					const scripts = document.querySelectorAll('script');
					for (const script of scripts) {
						const content = script.textContent;
						if (content && (content.includes('google.maps') || content.includes('LatLng'))) {
							const latMatch = content.match(/lat[:\s]*(-?\d+\.\d+)/i);
							const lngMatch = content.match(/lng[:\s]*(-?\d+\.\d+)/i);
							
							if (latMatch && latMatch[1]) {
								latitud = parseFloat(latMatch[1]);
							}
							
							if (lngMatch && lngMatch[1]) {
								longitud = parseFloat(lngMatch[1]);
							}
							
							if (latitud !== 0 && longitud !== 0) break;
						}
					}
				}

				// Método 4: Buscar en elementos con clase específica
				if (latitud === 0 && longitud === 0) {
					const mapClasses = ['.map', '.google-map', '.property-map', '.location-map', '#map'];
					for (const className of mapClasses) {
						try {
							const mapElement = document.querySelector(className);
							if (mapElement) {
								const lat = mapElement.getAttribute('data-lat') || mapElement.getAttribute('data-latitude');
								const lng = mapElement.getAttribute('data-lng') || mapElement.getAttribute('data-longitude');
								
								if (lat) latitud = parseFloat(lat);
								if (lng) longitud = parseFloat(lng);
								
								if (latitud !== 0 && longitud !== 0) break;
							}
						} catch (e) {
							console.error('Error al buscar mapa por clase:', e);
						}
					}
				}

				console.log('Coordenadas extraídas:', { latitud, longitud });

				// Extraer servicios, ambientes y adicionales
				function extractFeatures(title) {
					// Método 1: Buscar por título h2 o div.titulo2
					const sectionByH2 = Array.from(document.querySelectorAll('h2')).find(h => h.textContent.trim().toUpperCase() === title);
					const sectionByDiv = Array.from(document.querySelectorAll('div.titulo2')).find(div => div.textContent.trim().toUpperCase() === title);
					
					let features = [];
					
					// Si encontramos la sección por h2
					if (sectionByH2) {
						let element = sectionByH2.nextElementSibling;
						while (element && element.tagName !== 'H2') {
							if (element.querySelector('.fa-check') || element.querySelector('.fa-check-circle')) {
								const text = element.textContent.trim();
								if (text) features.push(text);
							}
							element = element.nextElementSibling;
						}
					}
					
					// Si encontramos la sección por div.titulo2
					if (sectionByDiv) {
						// Buscar la lista ul que sigue al título
						const ul = sectionByDiv.nextElementSibling;
						if (ul && ul.tagName === 'UL' && ul.classList.contains('ficha_ul')) {
							// Obtener todos los elementos li con ícono de check
							const items = ul.querySelectorAll('li');
							for (const item of items) {
								if (item.querySelector('.fa-check') || item.querySelector('.detalleColorC')) {
									const text = item.textContent.trim();
									if (text) features.push(text);
								}
							}
						}
					}
					
					// Método 2: Buscar directamente en listas con clase ficha_ul
					if (features.length === 0) {
						document.querySelectorAll('ul.ficha_ul').forEach(ul => {
							// Verificar si el título anterior es el que buscamos
							const prevElement = ul.previousElementSibling;
							if (prevElement && prevElement.textContent.trim().toUpperCase() === title) {
								const items = ul.querySelectorAll('li');
								for (const item of items) {
									if (item.querySelector('.fa-check') || item.querySelector('.detalleColorC')) {
										const text = item.textContent.trim();
										if (text) features.push(text);
									}
								}
							}
						});
					}
					
					// Método 3: Buscar en secciones con clase específica
					if (features.length === 0) {
						const sectionClass = title.toLowerCase();
						const sectionSelector = '.' + sectionClass;
						try {
							const section = document.querySelector(sectionSelector);
							if (section) {
								const items = section.querySelectorAll('li');
								for (const item of items) {
									if (item.querySelector('.fa-check') || item.querySelector('.detalleColorC')) {
										const text = item.textContent.trim();
										if (text) features.push(text);
									}
								}
							}
						} catch (e) {
							console.error('Error al buscar sección por clase:', e);
						}
					}
					
					// Método 4: Buscar en la descripción de la propiedad
					if (features.length === 0 && title === 'ADICIONALES') {
						const descripcionElement = document.querySelector('#prop-desc');
						if (descripcionElement) {
							const descripcion = descripcionElement.textContent.toLowerCase();
							
							// Lista de adicionales comunes para buscar en la descripción
							const adicionalesComunes = [
								'Calefacción', 'Apto profesional', 'Luminoso', 'Termo eléctrico',
								'Aire acondicionado', 'Balcón', 'Terraza', 'Parrilla', 'Piscina',
								'Gimnasio', 'Seguridad', 'Vigilancia', 'Portero eléctrico', 'Ascensor'
							];
							
							for (const adicional of adicionalesComunes) {
								if (descripcion.includes(adicional.toLowerCase())) {
									features.push(adicional);
								}
							}
						}
					}
					
					// Método 5: Buscar en la descripción de la propiedad para servicios
					if (features.length === 0 && title === 'SERVICIOS') {
						const descripcionElement = document.querySelector('#prop-desc');
						if (descripcionElement) {
							const descripcion = descripcionElement.textContent.toLowerCase();
							
							// Lista de servicios comunes para buscar en la descripción
							const serviciosComunes = [
								'Agua Corriente', 'Cloaca', 'Gas Natural', 'Electricidad',
								'Internet', 'Cable', 'Teléfono', 'Agua Caliente'
							];
							
							for (const servicio of serviciosComunes) {
								if (descripcion.includes(servicio.toLowerCase())) {
									features.push(servicio);
								}
							}
						}
					}
					
					// Método 6: Buscar en la descripción de la propiedad para ambientes
					if (features.length === 0 && title === 'AMBIENTES') {
						const descripcionElement = document.querySelector('#prop-desc');
						if (descripcionElement) {
							const descripcion = descripcionElement.textContent.toLowerCase();
							
							// Lista de ambientes comunes para buscar en la descripción
							const ambientesComunes = [
								'Cocina', 'Comedor', 'Living', 'Lavadero', 'Patio',
								'Jardín', 'Dormitorio', 'Baño', 'Vestidor', 'Estudio',
								'Oficina', 'Sala de estar', 'Hall', 'Recepción'
							];
							
							for (const ambiente of ambientesComunes) {
								if (descripcion.includes(ambiente.toLowerCase())) {
									features.push(ambiente);
								}
							}
						}
					}
					
					// Método 7: Buscar en elementos con clases específicas
					if (features.length === 0) {
						// Clases comunes que pueden contener características
						const classesToCheck = [
							'.caracteristicas', '.features', '.amenities', '.servicios',
							'.ambientes', '.adicionales', '.detalles', '.details'
						];
						
						for (const className of classesToCheck) {
							try {
								const elements = document.querySelectorAll(className + ' li');
								for (const element of elements) {
									const text = element.textContent.trim();
									if (text) features.push(text);
								}
							} catch (e) {
								console.error('Error al buscar por clase:', e);
							}
						}
					}
					
					// Eliminar duplicados
					features = [...new Set(features)];
					
					console.log('Características extraídas para ' + title + ':', features);
					return features;
				}
				
				let servicios = extractFeatures('SERVICIOS');
				let tiposAmbientes = extractFeatures('AMBIENTES');
				let adicionales = extractFeatures('ADICIONALES');

				console.log('Imágenes extraídas:', images); // Log para verificar imágenes
				console.log('Servicios extraídos:', servicios);
				console.log('Ambientes extraídos:', tiposAmbientes);
				console.log('Adicionales extraídos:', adicionales);

				// Extraer características de la descripción si no se encontraron de otra manera
				if (servicios.length === 0 && tiposAmbientes.length === 0 && adicionales.length === 0) {
					console.log('No se encontraron características en secciones específicas, intentando extraer de la descripción...');
					
					const descripcionElement = document.querySelector('#prop-desc');
					if (descripcionElement) {
						const descripcion = descripcionElement.textContent.toLowerCase();
						
						// Función para verificar si un término está en la descripción
						function containsTerm(text, term) {
							const termLower = term.toLowerCase();
							// Verificar el término exacto
							if (text.includes(termLower)) return true;
							
							// Verificar variaciones comunes
							const variations = [
								termLower,
								termLower.replace(/ó/g, 'o'),
								termLower.replace(/á/g, 'a'),
								termLower.replace(/é/g, 'e'),
								termLower.replace(/í/g, 'i'),
								termLower.replace(/ú/g, 'u'),
								termLower.replace(/ñ/g, 'n'),
								termLower + 's', // plural
								termLower.replace(/s$/, '') // singular
							];
							
							return variations.some(v => text.includes(v));
						}
						
						// Servicios comunes con términos alternativos
						const serviciosMap = {
							'Agua Corriente': ['agua', 'agua corriente'],
							'Cloaca': ['cloaca', 'cloacas', 'desagüe'],
							'Gas Natural': ['gas', 'gas natural', 'gas envasado'],
							'Electricidad': ['electricidad', 'luz', 'electrica'],
							'Internet': ['internet', 'wifi', 'wi-fi'],
							'Cable': ['cable', 'television', 'tv'],
							'Teléfono': ['telefono', 'linea telefonica'],
							'Agua Caliente': ['agua caliente', 'termotanque', 'calefon']
						};
						
						// Ambientes comunes con términos alternativos
						const ambientesMap = {
							'Cocina': ['cocina', 'kitchenette'],
							'Comedor': ['comedor', 'dining'],
							'Living': ['living', 'sala', 'estar'],
							'Lavadero': ['lavadero', 'lavanderia'],
							'Patio': ['patio', 'jardin', 'exterior'],
							'Jardín': ['jardin', 'verde', 'parque'],
							'Dormitorio': ['dormitorio', 'habitacion', 'cuarto', 'dorm'],
							'Baño': ['baño', 'toilette', 'sanitario'],
							'Vestidor': ['vestidor', 'walking closet', 'placard'],
							'Estudio': ['estudio', 'escritorio', 'home office'],
							'Oficina': ['oficina', 'despacho'],
							'Sala de estar': ['sala de estar', 'family room'],
							'Hall': ['hall', 'recibidor', 'entrada'],
							'Recepción': ['recepcion', 'lobby']
						};
						
						// Adicionales comunes con términos alternativos
						const adicionalesMap = {
							'Calefacción': ['calefaccion', 'calefactor', 'caldera', 'radiador'],
							'Apto profesional': ['apto profesional', 'uso profesional', 'consultorio'],
							'Luminoso': ['luminoso', 'luz natural', 'soleado', 'iluminado'],
							'Termo eléctrico': ['termo', 'termotanque', 'calefon'],
							'Aire acondicionado': ['aire', 'aire acondicionado', 'split', 'climatizacion'],
							'Balcón': ['balcon', 'terraza pequeña'],
							'Terraza': ['terraza', 'azotea', 'roof'],
							'Parrilla': ['parrilla', 'asador', 'barbacoa', 'bbq'],
							'Piscina': ['piscina', 'pileta', 'natatorio'],
							'Gimnasio': ['gimnasio', 'gym'],
							'Seguridad': ['seguridad', 'vigilancia', 'guardia'],
							'Vigilancia': ['vigilancia', 'seguridad 24hs', 'camaras'],
							'Portero eléctrico': ['portero', 'portero electrico', 'intercom'],
							'Ascensor': ['ascensor', 'elevador', 'lift']
						};
						
						// Extraer servicios
						for (const [servicio, terminos] of Object.entries(serviciosMap)) {
							if (terminos.some(term => containsTerm(descripcion, term))) {
								servicios.push(servicio);
							}
						}
						
						// Extraer ambientes
						for (const [ambiente, terminos] of Object.entries(ambientesMap)) {
							if (terminos.some(term => containsTerm(descripcion, term))) {
								tiposAmbientes.push(ambiente);
							}
						}
						
						// Extraer adicionales
						for (const [adicional, terminos] of Object.entries(adicionalesMap)) {
							if (terminos.some(term => containsTerm(descripcion, term))) {
								adicionales.push(adicional);
							}
						}
						
						// Eliminar duplicados
						servicios = [...new Set(servicios)];
						tiposAmbientes = [...new Set(tiposAmbientes)];
						adicionales = [...new Set(adicionales)];
						
						console.log('Características extraídas de la descripción:');
						console.log('- Servicios:', servicios);
						console.log('- Ambientes:', tiposAmbientes);
						console.log('- Adicionales:', adicionales);
					}
				}

				// Método 8: Buscar en tablas de detalles
				function extractFeaturesFromTables() {
					const tables = document.querySelectorAll('table');
					const featuresFound = {
						servicios: [],
						ambientes: [],
						adicionales: []
					};
					
					// Términos comunes para identificar categorías en tablas
					const categoryTerms = {
						servicios: ['servicios', 'services', 'utilities', 'instalaciones'],
						ambientes: ['ambientes', 'rooms', 'espacios', 'environments', 'distribución'],
						adicionales: ['adicionales', 'extras', 'amenities', 'comodidades', 'características']
					};
					
					// Términos para clasificar características
					const featureClassification = {
						servicios: ['agua', 'gas', 'luz', 'electricidad', 'cloaca', 'internet', 'cable', 'teléfono', 'telefono'],
						ambientes: ['cocina', 'comedor', 'living', 'dormitorio', 'habitacion', 'baño', 'lavadero', 'patio', 'jardin', 'estudio', 'oficina'],
						adicionales: ['calefaccion', 'aire', 'balcon', 'terraza', 'parrilla', 'piscina', 'pileta', 'gimnasio', 'seguridad', 'ascensor', 'portero', 'apto profesional']
					};
					
					// Procesar cada tabla
					tables.forEach(table => {
						// Verificar si la tabla tiene un título o encabezado que indique categoría
						let tableCategory = null;
						const tableCaption = table.querySelector('caption');
						const tableHeader = table.querySelector('th');
						
						if (tableCaption) {
							const captionText = tableCaption.textContent.toLowerCase();
							for (const [category, terms] of Object.entries(categoryTerms)) {
								if (terms.some(term => captionText.includes(term))) {
									tableCategory = category;
									break;
								}
							}
						} else if (tableHeader) {
							const headerText = tableHeader.textContent.toLowerCase();
							for (const [category, terms] of Object.entries(categoryTerms)) {
								if (terms.some(term => headerText.includes(term))) {
									tableCategory = category;
									break;
								}
							}
						}
						
						// Extraer todas las celdas de la tabla
						const cells = table.querySelectorAll('td');
						cells.forEach(cell => {
							const text = cell.textContent.trim();
							if (text) {
								// Si conocemos la categoría de la tabla, añadir a esa categoría
								if (tableCategory) {
									featuresFound[tableCategory].push(text);
								} else {
									// Intentar clasificar la característica
									const textLower = text.toLowerCase();
									let classified = false;
									
									for (const [category, terms] of Object.entries(featureClassification)) {
										if (terms.some(term => textLower.includes(term))) {
											featuresFound[category].push(text);
											classified = true;
											break;
										}
									}
									
									// Si no se pudo clasificar, añadir como adicional
									if (!classified && text.length > 2) {
										featuresFound.adicionales.push(text);
									}
								}
							}
						});
					});
					
					// Añadir las características encontradas a las listas principales
					if (servicios.length === 0 && featuresFound.servicios.length > 0) {
						servicios = featuresFound.servicios;
					}
					
					if (tiposAmbientes.length === 0 && featuresFound.ambientes.length > 0) {
						tiposAmbientes = featuresFound.ambientes;
					}
					
					if (adicionales.length === 0 && featuresFound.adicionales.length > 0) {
						adicionales = featuresFound.adicionales;
					}
					
					console.log('Características extraídas de tablas:', featuresFound);
				}
				
				// Ejecutar extracción de tablas si no se encontraron características
				if (servicios.length === 0 && tiposAmbientes.length === 0 && adicionales.length === 0) {
					extractFeaturesFromTables();
				}

				// Método 9: Buscar en elementos específicos de Tokko
				function extractFeaturesFromTokkoElements() {
					const featuresFound = {
						servicios: [],
						ambientes: [],
						adicionales: []
					};
					
					// Clases comunes en sitios Tokko
					const tokkoClasses = [
						'.ficha_detalle_item', '.ficha_detalle_datos', '.ficha_detalle_info',
						'.property-features', '.property-amenities', '.property-services',
						'.property-details', '.property-rooms', '.property-extras'
					];
					
					// Términos para clasificar características
					const featureClassification = {
						servicios: ['agua', 'gas', 'luz', 'electricidad', 'cloaca', 'internet', 'cable', 'teléfono', 'telefono'],
						ambientes: ['cocina', 'comedor', 'living', 'dormitorio', 'habitacion', 'baño', 'lavadero', 'patio', 'jardin', 'estudio', 'oficina'],
						adicionales: ['calefaccion', 'aire', 'balcon', 'terraza', 'parrilla', 'piscina', 'pileta', 'gimnasio', 'seguridad', 'ascensor', 'portero', 'apto profesional']
					};
					
					// Buscar en cada clase
					tokkoClasses.forEach(className => {
						try {
							const elements = document.querySelectorAll(className);
							elements.forEach(element => {
								// Buscar elementos li dentro del contenedor
								const items = element.querySelectorAll('li');
								if (items.length > 0) {
									items.forEach(item => {
										const text = item.textContent.trim();
										if (text) {
											// Intentar clasificar la característica
											const textLower = text.toLowerCase();
											let classified = false;
											
											for (const [category, terms] of Object.entries(featureClassification)) {
												if (terms.some(term => textLower.includes(term))) {
													featuresFound[category].push(text);
													classified = true;
													break;
												}
											}
											
											// Si no se pudo clasificar, añadir como adicional
											if (!classified && text.length > 2) {
												featuresFound.adicionales.push(text);
											}
										}
									});
								} else {
									// Si no hay elementos li, verificar el texto del elemento
									const text = element.textContent.trim();
									if (text) {
										// Buscar etiquetas y valores en el formato "Etiqueta: Valor"
										const colonIndex = text.indexOf(':');
										if (colonIndex !== -1) {
											const label = text.substring(0, colonIndex).trim().toLowerCase();
											const value = text.substring(colonIndex + 1).trim();
											
											if (value) {
												// Clasificar según la etiqueta
												if (label.includes('servicio') || label.includes('instalacion')) {
													featuresFound.servicios.push(value);
												} else if (label.includes('ambiente') || label.includes('habitacion') || label.includes('room')) {
													featuresFound.ambientes.push(value);
												} else if (label.includes('adicional') || label.includes('extra') || label.includes('amenity')) {
													featuresFound.adicionales.push(value);
												}
											}
										}
									}
								}
							});
						} catch (e) {
							console.error('Error al buscar en clase Tokko:', e);
						}
					});
					
					// Añadir las características encontradas a las listas principales
					if (servicios.length === 0 && featuresFound.servicios.length > 0) {
						servicios = featuresFound.servicios;
					}
					
					if (tiposAmbientes.length === 0 && featuresFound.ambientes.length > 0) {
						tiposAmbientes = featuresFound.ambientes;
					}
					
					if (adicionales.length === 0 && featuresFound.adicionales.length > 0) {
						adicionales = featuresFound.adicionales;
					}
					
					console.log('Características extraídas de elementos Tokko:', featuresFound);
				}
				
				// Ejecutar extracción de elementos Tokko si no se encontraron características
				if (servicios.length === 0 && tiposAmbientes.length === 0 && adicionales.length === 0) {
					extractFeaturesFromTokkoElements();
				}

				// Normalizar los nombres de las características
				function normalizeFeatures(features) {
					return features.map(feature => {
						// Eliminar espacios en blanco al inicio y final
						let normalizedFeature = feature.trim();
						
						// Capitalizar primera letra
						if (normalizedFeature.length > 0) {
							normalizedFeature = normalizedFeature.charAt(0).toUpperCase() + normalizedFeature.slice(1).toLowerCase();
						}
						
						return normalizedFeature;
					});
				}
				
				// Normalizar las características extraídas
				servicios = normalizeFeatures(servicios);
				tiposAmbientes = normalizeFeatures(tiposAmbientes);
				adicionales = normalizeFeatures(adicionales);

				console.log('Datos extraídos:', {
					tipoPropiedad, ubicacion, operacion, dormitorios, banios, antiguedad,
					superficieCubierta, superficieTotal, superficieTerreno,
					frente, fondo, ambientes, plantas, cocheras, situacion,
					expensas, descripcion, images, condicion, orientacion, disposicion,
					servicios, tiposAmbientes, adicionales
				});

				return {
					tipoPropiedad,
					ubicacion,
					operacion,
					dormitorios,
					banios,
					antiguedad,
					superficieCubierta,
					superficieTotal,
					superficieTerreno,
					frente,
					fondo,
					ambientes,
					plantas,
					cocheras,
					situacion,
					expensas,
					descripcion,
					images,
					condicion,
					orientacion,
					disposicion,
					servicios,
					tiposAmbientes,
					adicionales,
					latitud,
					longitud
				};
			})()
		`, &details),
	)

	if err != nil {
		return nil, fmt.Errorf("error extrayendo detalles: %v (url: %s)", err, url)
	}

	fmt.Printf("✓ Extracción completada: %+v\n", details)
	return &details, nil
}
