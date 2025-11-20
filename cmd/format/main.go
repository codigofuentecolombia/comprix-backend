package main

import (
	"comprix/app/utils"
	"fmt"
)

func main() {
	// Lista de pruebas con distintos formatos de moneda y su valor esperado
	testCases := []struct {
		input    string
		expected float64
	}{
		{"1338.83", 1338.83},
		{"$182.000", 182000},
		{"$1.200", 1200},
		{"$1,200", 1200},
		{"$1.000", 1000},
		{"1,1", 1.1},
		{"1.0", 1.0},
		{".9", 0.9},
		{",1", 0.1},
		{"10.12", 10.12},
		{"10,12", 10.12},
		{"2003.42", 2003.42},
		{"2,003.42", 2003.42},
		{"2.003,42", 2003.42},
		{"$2,003.42", 2003.42},
		{"€2.003,42", 2003.42},
		{"¥1,234,567.89", 1234567.89},
		{"1.234.567,89", 1234567.89},
		{"1,234,567.89", 1234567.89},
		{"1234567.89", 1234567.89},
		{"R$ 1.234,56", 1234.56},
		{"£1,234.56", 1234.56},
		{"asxasasxghkjn£1,234.56", 1234.56},
	}

	// Crear un escritor para las tablas con bordes
	printTableHeader()

	// Recorrer los casos de prueba y hacer los asserts
	for _, test := range testCases {
		cleanedValue, err := utils.CleanCurrencyFormat(test.input)
		if err != nil {
			// Imprimir fila con error
			printTableRow("Error", test.input, fmt.Sprintf("Error: %v", err))
		} else {
			// Verificar si el valor calculado es el esperado
			if cleanedValue != test.expected {
				// Imprimir fila con error en valor
				printTableRow("❌", test.input, fmt.Sprintf("%.2f", cleanedValue))
			} else {
				// Imprimir fila con éxito
				printTableRow("🟢", test.input, fmt.Sprintf("%.2f", cleanedValue))
			}
		}
	}

	printTableFooter()
}

func printTableHeader() {
	// Imprimir el encabezado con bordes
	fmt.Println("+------------+-------------------------+--------------------------------------+")
	fmt.Println("| Estado     | Valor Original          | Valor Convertido                     |")
	fmt.Println("+------------+-------------------------+--------------------------------------+")
}

func printTableRow(state, original, converted string) {
	// Imprimir las filas con bordes
	fmt.Printf("| %-10s | %-23s | %-34s |\n", state, original, converted)
}

func printTableFooter() {
	// Imprimir el borde final
	fmt.Println("+------------+-------------------------+--------------------------------------+")
}
