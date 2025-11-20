package pages

import (
	"comprix/app/domain/dto"
	"sync"
)

func (scrapper *ScrapPage) getCategoryPageLinks(totalPages int, url string) []string {
	var wg sync.WaitGroup
	// Crear canales
	linkChan := make(chan []string, totalPages)
	// Crear semaforos
	sem := make(chan struct{}, scrapper.maxGoRutines)
	// Iterar pagina por pagina
	for i := 1; i <= totalPages; i++ {
		wg.Add(1)
		// Bloquear si hay gorutinas pendientes
		sem <- struct{}{}
		// Go rutina
		go func(page int) {
			defer wg.Done()
			// Libera un espacio en el canal cuando finaliza
			defer func() { <-sem }()
			// Obtener links
			links := scrapper.svc.GetProductLinksByPage(dto.ScrapperParams{Url: url, Page: i})
			// Pasar links al canal
			linkChan <- links // Enviar los links al canal
		}(i)
	}
	// Gorutina para cerrar el canal cuando todas las gorutinas terminan
	go func() {
		wg.Wait()
		close(linkChan)
	}()
	// Recolectar los resultados del canal
	var links []string
	// Iterar canales
	for l := range linkChan {
		links = append(links, l...)
	}
	// Regresar valores
	return links
}

func (scrapper *ScrapPage) handleCategoryPageProducts(links []string, categories []string) {
	// Verificar si existen links
	if len(links) == 0 {
		return
	}
	// Usamos un WaitGroup para esperar a que todas las goroutines terminen
	var wg sync.WaitGroup
	// Limitar la concurrencia
	sem := make(chan struct{}, scrapper.maxGoRutines)
	// Lanzamos una goroutine para cada enlace
	for _, link := range links {
		// Incrementar el contador del WaitGroup
		wg.Add(1)
		// Goroutine para manejar el scraping de cada producto
		go func(url string) {
			// Decrementar el contador cuando la goroutine termine
			defer wg.Done()
			// Adquirir el semáforo (esto bloquea si hay más goroutines corriendo)
			sem <- struct{}{}
			// Crear o actualizar producto
			scrapper.CreateOrUpdateProduct(url, categories)
			// Liberar el semáforo (permitir que una nueva goroutine comience)
			<-sem
		}(link)
	}
	// Esperar a que todas las goroutines terminen
	wg.Wait()
}

func (scrapper *ScrapPage) ProcessCategoryLinks(links []dto.ScrapperCategoryLink, chunked ...bool) {
	var wg sync.WaitGroup
	// Especificar cuantos haran concurrencia
	jobs := make(chan dto.ScrapperCategoryLink, 1)
	// Iterar trabajos
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Iterar trabajos
			for link := range jobs {
				// Obtener el total de paginas
				totalPages := scrapper.svc.GetTotalPages(dto.ScrapperParams{Url: link.Link})
				// Obtener todos los links de los productos
				pageProductLinks := scrapper.getCategoryPageLinks(totalPages, link.Link)
				// Manejar la creacion o actualizacion de los productos
				if len(chunked) > 0 && chunked[0] {
					scrapper.processUrls(pageProductLinks, link.Categories)
				} else {
					scrapper.handleCategoryPageProducts(pageProductLinks, link.Categories)
					// Eliminar cache
					scrapper.removeChromeProcesses()
				}
			}
		}()
	}
	// Enviar trabajos
	for _, link := range links {
		jobs <- link
	}
	// Cerrar trabajos
	close(jobs)
	// Esperar a todos los workers
	wg.Wait()
	// Eliminar cache
	scrapper.removeChromeProcesses()
}
