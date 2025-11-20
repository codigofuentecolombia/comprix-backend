package service_product

import (
	"comprix/app/utils"
	"fmt"
	"strings"
)

func (svc *Service) SanitizeName(name string) string {
	name = strings.ToLower(name)           // Convertir a minusculas
	name = utils.NormalizeWhitespace(name) // Unificar espaciado
	name = utils.RemoveAccents(name)       // Quitar acentos
	// Formatear unicode
	name = svc.FormateUnicodeWords(name)
	// Eliminar caracteres especiales
	name = utils.RemoveSpecialChars(name)
	name = utils.SeparateNumbersAndWords(name)
	// Comprimir palabras
	name = fmt.Sprintf(" %s ", svc.ZipWords(name))
	// Remplazar vacio
	for _, word := range svc.GetEmptyReplacementWords() {
		name = strings.ReplaceAll(name, word, "")
	}
	// Remplazar por espacio
	for _, word := range svc.GetSpaceReplacementWords() {
		name = strings.ReplaceAll(name, word, " ")
	}
	// Ordenar palabras
	return utils.SortWords(name)
}
