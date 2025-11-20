package service_product

import (
	"comprix/app/domain/dto"
	"comprix/app/utils"
	"fmt"
	"slices"
	"strings"
)

func (svc *Service) SearchUnicode() {
	name := "�"
	list := []string{}
	// Obtener productos
	response := svc.repositories.product.FindAll(dto.ProductRepositoryFindParams{
		Name:    &name,
		Selects: &dto.RepositoryGormSelections{Query: []string{"id", "name"}},
	})
	// Iterar products
	for _, product := range response.Data {
		productName := utils.NormalizeWhitespace(strings.ToLower(product.Name))
		productName = utils.RemoveAccents(productName)
		productName = svc.FormateUnicodeWords(productName)
		// Separar por palabras
		for _, word := range strings.Split(productName, " ") {
			// Verificar si la palabra tiene el unicode
			if strings.Contains(word, name) && !slices.Contains(list, word) {
				list = append(list, word)
				fmt.Println(fmt.Sprintf("%s - %d", word, product.ID))
			}
		}
	}

}
