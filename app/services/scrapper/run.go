package service_scrapper

import (
	"comprix/app/utils"
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// Correr acciones
func (s *Service) RunActions(ctx context.Context, url string, actions []chromedp.Action) error {
	var readyState string
	// Acciones basicas
	baseActions := []chromedp.Action{
		// Bloquear recursos
		network.Enable(),
		network.SetBlockedURLS(utils.GetBlockedResources()),
		// Establecer dimensiones de escritorio
		chromedp.EmulateViewport(1500, 1000),
		// Navegar a la página
		chromedp.Navigate(url),
		// Esperar que el cuerpo de la página esté listo
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.WaitReady("body", chromedp.ByQuery),
		// Verificar el estado de carga de la página
		chromedp.Evaluate(`document.readyState`, &readyState),
		// Esperar hasta que el estado de carga sea "complete"
		chromedp.ActionFunc(func(ctx context.Context) error {
			total := 3
			current := 0
			// Espera hasta que readyState sea "complete"
			for readyState != `"complete"` && current < total {
				chromedp.Evaluate(`document.readyState`, &readyState)
				time.Sleep(500 * time.Millisecond)
				current++
			}
			return nil
		}),
		// Mover mouse
		chromedp.Evaluate(s.MouseMoveScript(100, 100, "document.body"), nil),
		chromedp.Evaluate(s.MouseMoveScript(100, 150, "document.body"), nil),
		chromedp.Evaluate(s.MouseMoveScript(100, 200, "document.body"), nil),
	}
	// Fusionar acciones base con acciones adicionales
	allActions := append(baseActions, actions...)
	// Correr
	return chromedp.Run(ctx, allActions...)
}

// Inicializar contexto y correr con acciones
func (s *Service) InitAndRunActions(url string, actions []chromedp.Action, optTimeout ...time.Duration) error {
	// Crear contexto
	ctx, cancelFns := s.InitContext(optTimeout...)
	// Asegurar cierre al salir
	defer cancelFns()
	// Correr acciones
	err := s.RunActions(ctx, url, actions)
	// Regresar error
	return err
}

func (s *Service) MouseMoveScript(coordX int, coordY int, element string) string {
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
