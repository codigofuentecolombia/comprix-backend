package utils

import "github.com/agnivade/levenshtein"

// Comparar que porcentaje de similitud hay entre strings
func LevenshteinSimilarity(a, b string) float64 {
	distance := levenshtein.ComputeDistance(a, b)
	maxLength := max(len(a), len(b))
	// Regresar porcentaje de similitud
	return 1 - float64(distance)/float64(maxLength)
}
