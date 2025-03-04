package scraper

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// Inmobiliaria representa una agencia inmobiliaria
type Inmobiliaria struct {
	Nombre    string
	SitioWeb  string
	Direccion string
	Telefono  string
	Rating    string
}

// Nombre del archivo CSV
const fileName = "inmobiliarias_lanus.csv"

// SearchInmobiliarias busca inmobiliarias en Google Maps y opcionalmente las guarda en CSV
func SearchInmobiliarias(ctx context.Context, zona string) ([]Inmobiliaria, error) {
	fmt.Println("🚀 Iniciando scraper...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1280, 800),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	url := fmt.Sprintf("https://www.google.com/maps/search/inmobiliarias+en+%s", strings.ReplaceAll(zona, " ", "+"))
	fmt.Printf("📍 Navegando a: %s\n", url)

	var results []Inmobiliaria

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(2*time.Second),
		chromedp.WaitVisible(`div[role="feed"]`),

		// Hacer scroll
		chromedp.Evaluate(`
			(() => {
				const feed = document.querySelector('div[role="feed"]');
				let lastHeight = feed.scrollHeight;
				let scrollCount = 0;
				const maxScrolls = 30;

				function doScroll() {
					feed.scrollTo(0, feed.scrollHeight);
					scrollCount++;
					
					setTimeout(() => {
						const newHeight = feed.scrollHeight;
						if (newHeight > lastHeight && scrollCount < maxScrolls) {
							lastHeight = newHeight;
							doScroll();
						}
					}, 2000);
				}

				doScroll();
			})()
		`, nil),

		chromedp.Sleep(30*time.Second),

		// Extraer información
		chromedp.Evaluate(`
			(() => {
				const results = [];
				const articles = Array.from(document.querySelectorAll('div.Nv2PK'));
				
				return articles.map(element => {
					const nombreElement = element.querySelector('div.qBF1Pd');
					const ratingElement = element.querySelector('span.MW4etd');
					const telefonoElement = element.querySelector('span.UsdlK');
					
					let sitioWeb = '';
					const links = Array.from(element.querySelectorAll('a[href*="http"]'));
					
					const realWebsite = links.find(a => 
						!a.href.includes('google.com/maps') && 
						!a.href.includes('search') &&
						a.href.includes('http')
					);

					if (realWebsite) {
						sitioWeb = realWebsite.href;
					}
					
					const direccionContainer = element.querySelector('div.W4Efsd');
					let direccion = '';
					if (direccionContainer) {
						const spans = Array.from(direccionContainer.querySelectorAll('span'));
						
						for (let i = 0; i < spans.length; i++) {
							const text = spans[i].textContent;
							if (text === 'Agencia inmobiliaria' || text === 'Agencia de bienes inmuebles comerciales') {
								for (let j = i + 1; j < i + 4 && j < spans.length; j++) {
									const nextText = spans[j].textContent;
									if (!nextText.includes('·') && 
										!nextText.includes('Agencia') &&
										!nextText.includes('Abierto') &&
										!nextText.includes('Cerrado') &&
										nextText.length > 5) {
										direccion = nextText.trim();
										break;
									}
								}
								break;
							}
						}

						if (!direccion) {
							const allW4Efsd = element.querySelectorAll('div.W4Efsd');
							if (allW4Efsd.length > 2) {
								const thirdW4Efsd = allW4Efsd[2];
								const thirdSpans = Array.from(thirdW4Efsd.querySelectorAll('span'));
								
								thirdSpans.forEach(span => {
									const text = span.textContent;
									if (!text.includes('·') && 
										!text.includes('Agencia') &&
										!text.includes('Abierto') &&
										!text.includes('Cerrado') &&
										text.length > 5) {
										direccion = text.trim();
									}
								});
							}
						}
					}
					
					return {
						nombre: nombreElement ? nombreElement.innerText.trim() : '',
						rating: ratingElement ? ratingElement.textContent.trim() : '',
						direccion: direccion,
						telefono: telefonoElement ? telefonoElement.innerText.trim() : '',
						sitioWeb: sitioWeb
					};
				});
			})()
		`, &results),
	)

	if err != nil {
		return nil, fmt.Errorf("error ejecutando script: %v", err)
	}

	fmt.Printf("🔍 Encontradas %d inmobiliarias\n", len(results))
	return results, nil
}

func saveToCSV(results []Inmobiliaria, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir headers
	headers := []string{"Nombre", "Sitio Web", "Dirección", "Teléfono", "Rating"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Escribir resultados
	for _, r := range results {
		row := []string{
			r.Nombre,
			r.SitioWeb,
			r.Direccion,
			r.Telefono,
			r.Rating,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// limpiarTexto elimina caracteres no deseados y espacios extra
func limpiarTexto(texto string) string {
	// Eliminar saltos de línea y espacios múltiples
	texto = strings.ReplaceAll(texto, "\n", " ")
	texto = strings.Join(strings.Fields(texto), " ")
	return strings.TrimSpace(texto)
}

// Agregar estas funciones de limpieza
func limpiarDireccion(direccion string) string {
	direccion = strings.TrimSpace(direccion)

	// Eliminar el punto y espacios al principio
	direccion = strings.TrimPrefix(direccion, "·")
	direccion = strings.TrimPrefix(direccion, " ")

	// Si es un teléfono o texto de horario, ignorarlo
	if strings.HasPrefix(direccion, "011") ||
		strings.Contains(direccion, "Abre") ||
		strings.Contains(direccion, "Cierra") ||
		strings.Contains(direccion, "⋅") {
		return ""
	}

	return direccion
}

func limpiarTelefono(telefono string) string {
	// Asegurar formato consistente de teléfono
	telefono = strings.TrimSpace(telefono)
	if telefono != "" && !strings.HasPrefix(telefono, "011") {
		telefono = "011 " + telefono
	}
	return telefono
}

func limpiarSitioWeb(url string) string {
	// Eliminar URLs de Google Ads y parámetros de tracking
	if strings.Contains(url, "googleadservices.com") ||
		strings.Contains(url, "?utm_source=") {
		return ""
	}
	return strings.TrimSpace(url)
}
