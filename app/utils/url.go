package utils

import (
	"net/url"
)

func CheckIfUrlHasQueryParams(link string) bool {
	parsedURL, err := url.Parse(link)
	// Si no se puede parsear, asumimos que no tiene query params
	if err != nil {
		return false
	}
	// Regresar
	return len(parsedURL.Query()) > 0
}
