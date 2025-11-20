package service_scrapper

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// Espera a que un elemento esté visible. Si no aparece tras varios intentos, toma un screenshot.
func (s *Service) WaitForElementOrScreenshot(selector string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		total := 3
		current := 0
		visible := false
		// Espera hasta que readyState sea "complete"
		for !visible && current < total {
			fmt.Printf("Vuelta: #%v\n", current)
			chromedp.Run(ctx, chromedp.Evaluate(fmt.Sprintf(`!!document.querySelector("%s")`, selector), &visible))
			time.Sleep(10000 * time.Millisecond)
			current++
		}
		// Verificar si es visible
		var screenshot []byte
		// Tomar captura
		fmt.Println(chromedp.Run(ctx, chromedp.CaptureScreenshot(&screenshot)))
		// Verificar si hay error
		os.WriteFile("error_screenshot.png", screenshot, 0644)

		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		// Regresar sin error
		return nil
	}
}

func (s *Service) GetSecuredValue(selector string, result *string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var exists bool
		// Verificar si el elemento existe en la página
		err := chromedp.Run(ctx,
			chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s") !== null`, selector), &exists),
		)
		// Validar si hay error
		if err != nil {
			return err
		} else {
			// Si el elemento existe, obtener el texto
			if exists {
				return chromedp.Text(selector, result, chromedp.ByQuery).Do(ctx)
			}
		}
		// Regresar sin error
		return nil
	}
}
