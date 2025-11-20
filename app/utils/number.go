package utils

func GetMinNumber(numbers []int) int {
	min := numbers[0]
	// Iterar numeros
	for _, num := range numbers {
		if num < min {
			min = num
		}
	}
	// Regresar minimo
	return min
}
