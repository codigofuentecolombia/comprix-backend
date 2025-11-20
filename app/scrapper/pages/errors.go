package pages

import (
	"comprix/app/domain/dao"
	"sync"
)

func (scrapper *ScrapPage) FixErrors() {
	var wg sync.WaitGroup
	// Especificar cuantos haran concurrencia
	jobs := make(chan dao.Error, scrapper.maxGoRutines)
	// Iterar trabajos
	for i := 0; i < scrapper.maxGoRutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Iterar trabajos
			for err := range jobs {
				// Procesar solo cuando haya categorias
				if len(err.Categories) > 0 {
					scrapper.CreateOrUpdateProduct(err.Url, err.Categories)
				}
			}
		}()
	}
	// Enviar trabajos
	for _, err := range scrapper.svc.GetErrors() {
		jobs <- err
	}
	// Cerrar trabajos
	close(jobs)
	// Esperar a todos los workers
	wg.Wait()
	// Eliminar cache
	scrapper.removeChromeProcesses()
}
