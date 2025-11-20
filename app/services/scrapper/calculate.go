package service_scrapper

import (
	"comprix/app/fails"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func (s *Service) CalculateTotalPages(totalProducts string) (int, error) {
	// Separar los componentes del texto
	parts := strings.Fields(totalProducts)
	// Verificar si hubo error
	if len(parts) == 0 {
		return 0, fails.Create("Utils CalculateTotalPages(): El texto esta vacio", nil)
	}
	// Intentar convertir el primer elemento a un número
	number, err := strconv.Atoi(parts[0])
	// Verificar si hubo error
	if err != nil {
		return 0, fails.Create(fmt.Sprintf("Utils CalculateTotalPages(): Error al convertir a número '%s'", totalProducts), err)
	}
	// Dividir por total de productos por pagina
	totalPages := int(math.Ceil(float64(number) / s.ProductsPerPage))
	// Verificar que no este vacio
	if totalPages < 1 {
		totalPages = 1
	}
	// Regresar
	return totalPages, nil
}
