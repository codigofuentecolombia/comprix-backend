package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

func NormalizeWhitespace(val string) string {
	// Formalizar espacios
	val = regexp.MustCompile(`\s+`).ReplaceAllString(val, " ")
	// Sustituir guion por espacio
	val = strings.ReplaceAll(val, "-", " ")
	// Sustituir coma por espacio
	val = NormalizeWhitespaceOnNumeric(strings.Split(val, ","))
	// Sustituir punto por espacio
	return NormalizeWhitespaceOnNumeric(strings.Split(val, "."))
}

func ExtractLetters(s string) string {
	return strings.Join(regexp.MustCompile(`[a-zA-Z]+`).FindAllString(s, -1), "")
}

func ExtractNumbers(s string) string {
	return strings.Join(regexp.MustCompile(`\d+`).FindAllString(s, -1), "")
}

func RemoveSpecialChars(s string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9\s]`).ReplaceAllString(s, "")
}

func SeparateNumbersAndWords(s string) string {
	return regexp.MustCompile(`([a-zA-Z]+)(\d+)|(\d+)([a-zA-Z]+)`).ReplaceAllString(s, `$1$3 $2$4 `)
}

func NormalizeWhitespaceOnNumeric(formatedVal []string) string {
	// Variable de retorno
	values := []string{}
	// Iterar
	if len(formatedVal) > 0 {
		for index := 0; index < len(formatedVal); index++ {
			word := formatedVal[index]
			// Verificar si hay elementos antes y después
			if index > 0 {
				lastWord := formatedVal[index-1]
				// Validar si el último carácter del anterior es dígito
				if len(lastWord) > 0 && len(word) > 0 &&
					unicode.IsDigit(rune(lastWord[len(lastWord)-1])) &&
					unicode.IsDigit(rune(word[0])) {
					// Agregar caracter
					values[len(values)-1] += fmt.Sprintf("0%s", word)
					continue
				}
			}
			// agregar al listado
			values = append(values, word)
		}
	}
	// Regresar
	return strings.Join(values, "")
}
