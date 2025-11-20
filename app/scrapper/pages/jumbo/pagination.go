package pages_jumbo

import (
	"comprix/app/domain/dto"
	service_scrapper "comprix/app/services/scrapper"
	"fmt"

	"github.com/chromedp/chromedp"
)

func (service Service) GetTotalPages(params dto.ScrapperParams) int {
	service.svc.CreateDebugLog(fmt.Sprintf("Iniciando scraping de la categoria. %s", params.Url), nil)
	// Obtener detalle de la pagina
	var totalPages int
	// Consultar pagina
	err := service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Esperar a que el contenido este cargado
		chromedp.WaitVisible(`#gallery-layout-container`, chromedp.ByQuery),
		// Obtener total de productos
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetTotalPagesScript()), &totalPages),
	}, 45)
	//Verificar si existe un error
	if err != nil {
		return service_scrapper.HandleTotalTries(service.svc, dto.HandleScrapperTotalTriesParams[int]{
			Err:             err,
			Url:             params.Url,
			Msg:             "No se pudo obtener el total de paginas",
			Callback:        service.GetTotalPages,
			CallbackArgs:    params,
			DefaultResponse: 1,
		})
	}
	// Regresar links
	return totalPages
}

func (service Service) GetProductLinksByPage(params dto.ScrapperParams) []string {
	// Variables
	var links []string
	var dummy interface{}
	// Definir pagina actual
	currentPageUrl := fmt.Sprintf("%s?page=%d", params.Url, params.Page)
	// Guaradr log
	service.svc.CreateDebugLog(fmt.Sprintf("Obteniendo detalle de pagina %s", currentPageUrl), nil)
	// Obtener detalles
	err := service.svc.InitAndRunActions(currentPageUrl, []chromedp.Action{
		// Esperar hasta que este listo
		chromedp.WaitVisible(`#gallery-layout-container`, chromedp.ByID),
		chromedp.Click(`#gallery-layout-container`, chromedp.ByID),
		// Recargar todos los productos
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.LoadAllPageProductsScript()), &dummy),
		// Esperar hasta que este listo
		chromedp.WaitVisible(`#se-ha-completado`, chromedp.ByID),
		// Obtener links
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetPageProductLinksScript()), &links),
	}, 45)
	// Validar el error
	if err != nil {
		return service_scrapper.HandleTotalTries(service.svc, dto.HandleScrapperTotalTriesParams[[]string]{
			Err:             err,
			Url:             currentPageUrl,
			Msg:             "No se pudo obtener los links de los productos",
			Callback:        service.GetProductLinksByPage,
			CallbackArgs:    params,
			DefaultResponse: []string{},
		})
	}
	// Regresar data
	return links
}
