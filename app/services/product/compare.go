package service_product

import (
	"comprix/app/utils"
	"fmt"
)

func (svc *Service) GetComparationPercentage(a, b string) float64 {
	// Obtener numeros
	numberA := utils.ExtractNumbers(a)
	numberB := utils.ExtractNumbers(b)
	// Verificar si no hay numeros
	hasNoNumbers := (numberA == "" && numberB == "")
	// Comparar el porcentaje de coincidencia
	if (utils.LevenshteinSimilarity(numberA, numberB) == 1) || hasNoNumbers {
		return utils.LevenshteinSimilarity(utils.ExtractLetters(a), utils.ExtractLetters(b))
	}
	// Regresar 0
	return 0
}

func (svc *Service) CompareNames(a, b string, shouldPrint ...bool) bool {
	percentage := svc.GetComparationPercentage(a, b)
	// Validar si se imprimira
	if len(shouldPrint) > 0 && shouldPrint[0] {
		fmt.Printf("Nombre A: %s\n", a)
		fmt.Printf("Nombre B: %s\n", b)
		fmt.Printf("Porcentaje: %f\n", percentage)
	}
	// Verificar si cumple con el minimo
	return percentage >= 0.875
}

func (svc *Service) CompareNamesWithDetails(a string, b string) bool {
	return svc.CompareNames(svc.SanitizeName(a), svc.SanitizeName(b), true)
}
