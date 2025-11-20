package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func CleanCurrencyFormat(input string) (float64, error) {
	// Eliminar caracteres no numéricos (excepto `.` y `,`)
	re := regexp.MustCompile(`[^\d.,]`)
	cleaned := AddDecimalsIfNeeded(re.ReplaceAllString(input, ""))
	// Detectar qué carácter es el separador decimal (último `,` o `.` en la cadena)
	lastComma := strings.LastIndex(cleaned, ",")
	lastDot := strings.LastIndex(cleaned, ".")
	//
	var decimalSeparator, thousandSeparator string
	// Validar posicion
	if lastComma > lastDot {
		decimalSeparator = ","
		thousandSeparator = "."
	} else {
		decimalSeparator = "."
		thousandSeparator = ","
	}
	// Eliminar correctamente los separadores de miles
	cleaned = strings.ReplaceAll(cleaned, thousandSeparator, "")
	// Reemplazar el separador decimal por `.`
	cleaned = strings.Replace(cleaned, decimalSeparator, ".", 1)
	// Convertir a float
	value, err := strconv.ParseFloat(cleaned, 64)
	if err != nil {
		return 0, err
	}
	// Regresar valor
	return value, nil
}

func AddDecimalsIfNeeded(input string) string {
	// Validar que al menos tenga 4 carateres
	if len(input) < 4 {
		return input
	}
	// Obtener el indici
	digitChar := input[len(input)-4]
	// Verificar si el cuarto carácter desde el final es un punto o coma
	if digitChar == '.' {
		// Verificar si hay un punto o coma en los últimos 3 caracteres
		if !strings.ContainsAny(input[len(input)-3:], ".,") {
			return input + ",00"
		}
	} else if digitChar == ',' {
		// Verificar si hay un punto o coma en los últimos 3 caracteres
		if !strings.ContainsAny(input[len(input)-3:], ".,") {
			return input + ".00"
		}
	}
	// Regresar igual
	return input
}
