package pages_carrefour

import (
	"comprix/app/domain/dto"
	"comprix/app/utils"
	"fmt"

	"github.com/chromedp/chromedp"

	service_scrapper "comprix/app/services/scrapper"
)

func (service Service) GetTotalPages(params dto.ScrapperParams) int {
	service.svc.CreateDebugLog(fmt.Sprintf("Iniciando scraping de la categoria. %s", params.Url), nil)
	// Obtener detalle de la pagina
	var totalPages int
	// Consultar pagina
	err := service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Esperar a que el contenido este cargado
		chromedp.WaitVisible(`.valtech-carrefourar-search-result-3-x-totalProducts--layout`, chromedp.ByQuery),
		// Obtener total de productos
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetTotalPagesScript()), &totalPages),
	}, 80)
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
	// Verificar si es mayor a cero
	if totalPages < 1 {
		return 1
	}
	// Regresar total de paginas
	return totalPages
}

func (service Service) GetProductLinksByPage(params dto.ScrapperParams) []string {
	// Variables
	var links []string
	var dummy interface{}
	var currentPageUrl string
	// Definir pagina actual
	if utils.CheckIfUrlHasQueryParams(params.Url) {
		currentPageUrl = fmt.Sprintf("%s&page=%d", params.Url, params.Page)
	} else {
		currentPageUrl = fmt.Sprintf("%s?page=%d", params.Url, params.Page)
	}
	// Guaradr log
	service.svc.CreateDebugLog(fmt.Sprintf("Obteniendo detalle de pagina %s", currentPageUrl), nil)
	// Obtener detalles
	err := service.svc.InitAndRunActions(currentPageUrl, []chromedp.Action{
		// Esperar a que el contenido este cargado
		chromedp.WaitVisible(`.valtech-carrefourar-search-result-3-x-gallery`, chromedp.ByQuery),
		// Simular mouse
		chromedp.ActionFunc(service.svc.MouseActions(".vtex-store-footer-2-x-footerLayout")),
		// Recargar todos los productos
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.LoadAllPageProducts()), &dummy),
		// Esperar hasta que este listo
		chromedp.WaitVisible(`#se-ha-completado`, chromedp.ByID),
		// Obtener total de productos
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetPageProductLinksScript()), &links),
	}, 150)
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
