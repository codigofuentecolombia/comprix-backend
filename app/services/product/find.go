package service_product

import (
	"comprix/app/domain/dao"
	"sync"
)

func (svc *Service) findMatches(product dao.Product, products []dao.Product, idsInUse map[uint]bool) (int, []string) {
	var mu sync.Mutex
	var wg sync.WaitGroup
	// Definir variables iniciales
	count := int(0)
	workers := 10
	matches := []string{}
	// Especificar cuantos haran concurrencia
	jobs := make(chan dao.Product, workers)
	// Iterar trabajos
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Iterar trabajos
			for compareProduct := range jobs {
				// Verificar si ya se procesó
				mu.Lock()
				if _, ok := idsInUse[compareProduct.ID]; ok {
					mu.Unlock()
					continue
				}
				mu.Unlock()
				// Verificar si cumple la condición
				if svc.CompareNames(product.Name, compareProduct.Name) {
					mu.Lock()
					if _, ok := idsInUse[compareProduct.ID]; !ok {
						count++
						matches = append(matches, compareProduct.Name)
						idsInUse[compareProduct.ID] = true
					}
					mu.Unlock()
				}
			}
		}()
	}
	// Enviar trabajos
	for subindex := 5451; subindex < len(products); subindex++ {
		jobs <- products[subindex]
	}
	// Cerrar trabajos
	close(jobs)
	// Esperar a todos los workers
	wg.Wait()
	// Regresar resultados
	return count, matches
}
