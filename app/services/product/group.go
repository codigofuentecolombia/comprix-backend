package service_product

import (
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"encoding/json"
	"fmt"
	"os"
)

func (svc *Service) AggroupByNames() error {
	// Obtener productos
	response := svc.repositories.product.FindAll(dto.ProductRepositoryFindParams{
		Selects: &dto.RepositoryGormSelections{Query: []string{"id", "name"}},
	})
	// Validar si hubo error
	if response.Error != nil {
		return fails.Create("GroupByNames", response.Error)
	}
	// Limpiar nombres
	for i := range response.Data {
		fmt.Printf("Sanitizando nombre de producto #%d \n", response.Data[i].ID)
		response.Data[i].Description = response.Data[i].Name
		response.Data[i].Name = svc.SanitizeName(response.Data[i].Name)
	}
	// Declarar variables
	count := 0
	idsInUse := map[uint]bool{}
	totalProducts := len(response.Data)
	groupedProducts := map[string][]string{}
	// Iterar productos
	for index := 0; index < totalProducts; index++ {
		// Agregar producto actual a los id en uso
		product := response.Data[index]
		idsInUse[product.ID] = true
		// Iterar productos para determinar cuales estan duplicados
		counter, matches := svc.findMatches(product, response.Data, idsInUse)
		// Actualizar variables solo si se encontraron registros
		if len(matches) > 0 {
			count += counter
			count++
			groupedProducts[product.Name] = matches
		}
		// Mostrar log
		fmt.Printf("Index %d de %d\n", (index + 1), totalProducts)
	}
	// Guardar en archivo
	if err := svc.SaveGroupingsToFile(groupedProducts); err != nil {
		return err
	}
	//
	fmt.Printf("Se han encontrado #%d - %d\n", count, len(idsInUse))
	//
	return nil
}

// Guardar agrupaciones en archivo
func (svc *Service) SaveGroupingsToFile(data map[string][]string) error {
	// Convertir a JSON con indentación
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	// Verificar si ocurrio un error al guardar el archivo
	if err != nil {
		return fails.Create("SaveGroupingsToFile: could not parse data", err)
	}
	// Crear archivo
	file, err := os.Create("datos.json")
	// Verificar si se pudo crear
	if err != nil {
		return fails.Create("SaveGroupingsToFile: could not create file", err)
	}
	defer file.Close()
	// Escribir al archivo
	_, err = file.Write(jsonBytes)
	// Verificar si se pudo escribir
	if err != nil {
		return fails.Create("SaveGroupingsToFile: could not write data", err)
	}
	// Regresar sin error
	return nil
}
