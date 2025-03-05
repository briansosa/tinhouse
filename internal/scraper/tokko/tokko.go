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

				console.log('Imágenes extraídas:', images); // Log para verificar imágenes

				console.log('Datos extraídos:', {
					tipoPropiedad, ubicacion, operacion, dormitorios, banios, antiguedad,
					superficieCubierta, superficieTotal, superficieTerreno,
					frente, fondo, ambientes, plantas, cocheras, situacion,
					expensas, descripcion, images
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
					disposicion
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
