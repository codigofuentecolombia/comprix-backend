package service_scrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

func (s *Service) MouseActions(selector string) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		var exists bool
		// Verificar si el elemento existe
		maxRetries := 3                         // Número máximo de intentos
		retryInterval := 500 * time.Millisecond // Tiempo de espera entre intentos
		// Intentar verificar la existencia del elemento con reintentos
		for i := 0; i < maxRetries; i++ {
			chromedp.Run(ctx,
				chromedp.Evaluate(fmt.Sprintf(`document.querySelector("%s") !== null`, selector), &exists),
			)
			// Si el elemento existe, salir del bucle
			if exists {
				break
			}
			// Esperar antes del siguiente intento
			time.Sleep(retryInterval)
		}
		// Verificar si existe
		if exists {
			// Emular movimiento del mouse y clic
			return chromedp.Run(ctx,
				// Mueve el mouse al elemento
				chromedp.Evaluate(s.MouseMoveAction(100, 100, fmt.Sprintf(`document.querySelector("%s")`, selector)), nil),
				// Espera un poco para que parezca más humano
				chromedp.Sleep(500*time.Millisecond),
				// Mover mouse
				chromedp.Evaluate(s.MouseMoveAction(100, 150, fmt.Sprintf(`document.querySelector("%s")`, selector)), nil),
				// Espera un poco para que parezca más humano
				chromedp.Sleep(500*time.Millisecond),
				// Mover mouse
				chromedp.Evaluate(s.MouseMoveAction(100, 200, fmt.Sprintf(`document.querySelector("%s")`, selector)), nil),
				// Espera un poco para que parezca más humano
				chromedp.Sleep(500*time.Millisecond),
				// Hace clic en el elemento
				chromedp.Click(selector, chromedp.ByQuery),
			)
		}
		// No hay error, solo que el elemento no existe
		return nil
	}
}

func (s *Service) MouseMoveAction(coordX int, coordY int, element string) string {
	return fmt.Sprintf(`
		if( %s ){
			%s.dispatchEvent(
				new MouseEvent('mousemove', {
					clientX: %d,
					clientY: %d,
					bubbles: true
				})
			);
			// Dar click
			%s.click();
		}
	`, element, element, coordX, coordY, element)
}
