package main

import (
	"fmt"
	"log"
	"slices"

	"github.com/xuri/excelize/v2"
)

type Category struct {
	Product     string   `json:"product"`
	Subcategory string   `json:"subcategory"`
	Category    string   `json:"category"`
	Links       []string `json:"links"`
}

func main() {
	filePath := "test.xlsx"
	// Abrimos el archivo Excel
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)

	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatalf("Error al leer filas: %v", err)
	}

	var categories []Category
	var proccessedLinks []string

	// Asumimos que la primera fila es encabezado
	for i, row := range rows {
		if i == 0 {
			continue // saltamos encabezados
		}
		// Validamos que haya al menos 4 columnas
		if len(row) < 4 {
			fmt.Printf("Se omitio el index %d\n", i+1)
			continue
		}
		// Verificar si link no esta vacio
		if row[3] != "" {
			var links []string
			currentIndex := 3
			// Verificar cantidad
			for {
				if currentIndex >= len(row) || row[currentIndex] == "" {
					break
				}
				// Obtener link
				link := row[currentIndex]
				// Verificar si ya esta en uso
				if !slices.Contains(proccessedLinks, link) {
					links = append(links, link)
					proccessedLinks = append(proccessedLinks, link)
				}
				// Incrementar contador
				currentIndex++
			}
			// Verificar si hay links
			if len(links) > 0 {
				// Agregar categoria
				categories = append(categories, Category{
					Links:       links,
					Product:     row[0],
					Category:    row[2],
					Subcategory: row[1],
				})
			}
		}
	}
	// Iterar categorias para imprimir
	for _, category := range categories {
		// Iterar links
		for _, link := range category.Links {
			fmt.Printf(`{Link: "%s", Categories: []string{"%s", "%s", "%s"}},`, link, category.Category, category.Subcategory, category.Product)
			fmt.Println("")
		}
	}
}
