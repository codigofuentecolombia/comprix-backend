package utils

import (
	"regexp"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func RemoveDuplicates(values []string) []string {
	// Usar un mapa para almacenar los elementos únicos
	unique := make(map[string]struct{})
	// Iterar sobre el slice y agregar los elementos al mapa
	for _, str := range values {
		unique[str] = struct{}{}
	}
	// Crear un slice para los elementos únicos
	result := make([]string, 0, len(unique))
	for str := range unique {
		result = append(result, str)
	}
	// Regresar resultado
	return result
}

func RemoveAllSpaces(value string) string {
	return strings.ReplaceAll(value, " ", "")
}

func SanitizeString(name string) string {
	// Convertir a minúsculas
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "0", "o")
	// Eliminar caracteres especiales
	re := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	name = re.ReplaceAllString(name, "")
	// Remplazar comunes
	name = strings.ReplaceAll(name, " del ", " ")
	name = strings.ReplaceAll(name, " las ", " ")
	name = strings.ReplaceAll(name, " los ", " ")
	nameWords := []string{}
	// Iterar por palabras
	if len(name) >= 6 {
		for _, word := range strings.Split(name, " ") {
			// verificar la longitud
			if len(word) > 2 {
				nameWords = append(nameWords, word)
			}
		}
		// Actualizar name
		name = strings.Join(nameWords, " ")
	}
	// Verificar tamaño
	if len(name) <= 4 {
		name = strings.ReplaceAll(name, " ", "")
	}
	// Dividir en palabras y ordenar
	words := strings.Fields(name)
	sort.Strings(words)
	// Volver a unir las palabras ordenadas
	return strings.Join(words, " ")
}

func RemoveAccents(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	// Transformar string
	result, _, _ := transform.String(t, input)
	// Regresar string sin acentos
	return result
}

func SortWords(input string) string {
	// Dividir en palabras y ordenar
	words := RemoveDuplicates(strings.Fields(input))
	sort.Strings(words)
	// Volver a unir las palabras ordenadas
	return strings.Join(words, " ")
}

func CheckIfStringIsNotEmpty(value *string) bool {
	return value != nil && *value != ""
}
