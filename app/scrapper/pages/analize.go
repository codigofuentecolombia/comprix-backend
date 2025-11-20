package pages

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/utils"
	"fmt"
	"sync"
	"time"
)

func AnalizeExistingPageProducts(conf *dto.Config, page dto.InitScrapPageFn, maxGoRutines int) {
	scrapper, err := InitService(conf, page, maxGoRutines)
	// Verificar si hubo error para finalizar funcion
	if err != nil {
		return
	}
	// Desactivar productos
	scrapper.disableProduct = true
	// Obtener productos
	products, err := scrapper.svc.GetExistingProducts()
	// Validar si se pudo obtener el producto
	if err != nil {
		scrapper.log.Error(fmt.Sprintf("No se pudo obtener los productos existentes de la pagina: error %v", err))
		// Finalizar
		return
	}
	// Obtener links de las categorias
	for index, chunkedProducts := range utils.ChunkSlice(products, 50) {
		start := time.Now()
		// Manejar la creacion o actualizacion de los productos
		scrapper.handleCategoryPageProducts(scrapper.ExtractUrlsFromPageProducts(chunkedProducts), []string{})
		// Eliminar cache
		scrapper.removeChromeProcesses()
		// Fecha fin
		end := time.Now()
		// Esperar 2 minutos
		time.Sleep(2 * time.Minute)
		// Mostrar mensaje
		scrapper.log.Debug(fmt.Sprintf("Se ha terminado de procesar el batch #%d, con una duración de %s", (index + 1), utils.FormatDuration(start, end)))
	}
}

func (svc *ScrapPage) ExtractUrlsFromPageProducts(products []dao.PageProduct) []string {
	// Canal para almacenar URLs
	urls := make(chan string, len(products))
	// WaitGroup para sincronizar goroutines
	var wg sync.WaitGroup
	// Iniciar goroutines para procesar las entidades
	for _, entity := range products {
		wg.Add(1)
		go func(e dao.PageProduct) {
			defer wg.Done()
			// Enviar URL al canal
			urls <- e.Url
		}(entity)
	}
	// Goroutine para cerrar el canal cuando terminen todas las goroutines
	go func() {
		wg.Wait()
		close(urls)
	}()
	// Convertir el canal en un slice de strings
	var urlList []string
	// Agregar url a la lista
	for url := range urls {
		urlList = append(urlList, url)
	}
	// Regresar
	return urlList
}
