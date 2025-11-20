package pages_masonline

import (
	"comprix/app/domain/dto"
	service_scrapper "comprix/app/services/scrapper"
	"comprix/app/utils"
	"fmt"

	"github.com/chromedp/chromedp"
)

func (service Service) GetTotalPages(params dto.ScrapperParams) int {
	service.svc.CreateDebugLog(fmt.Sprintf("Iniciando scraping de la categoria. %s", params.Url), nil)
	// Obtener detalle de la pagina
	var totalProducts string
	// Consultar pagina
	err := service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Esperar hasta que este listo
		chromedp.WaitVisible(".valtech-gdn-search-result-0-x-totalProducts--layout > span", chromedp.ByQuery),
		chromedp.Text(`.valtech-gdn-search-result-0-x-totalProducts--layout > span`, &totalProducts, chromedp.NodeVisible),
	}, 80)
	// Obtener total de productos
	totalPages, totalPagesErr := service.svc.CalculateTotalPages(totalProducts)
	//Verificar si existe un error
	if err != nil || totalPagesErr != nil {
		// Verificar el error
		if err == nil {
			err = totalPagesErr
		}
		// Regresar
		return service_scrapper.HandleTotalTries(service.svc, dto.HandleScrapperTotalTriesParams[int]{
			Err:             err,
			Url:             params.Url,
			Msg:             "No se pudo obtener el total de paginas",
			Callback:        service.GetTotalPages,
			CallbackArgs:    params,
			DefaultResponse: 1,
		})
	}
	// Verificar que no pase de 50
	if totalPages > 50 {
		totalPages = 50
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
		chromedp.WaitVisible(".valtech-gdn-search-result-0-x-totalProducts--layout span", chromedp.ByQuery),
		chromedp.WaitVisible(`.valtech-gdn-search-result-0-x-gallery`, chromedp.ByQuery),
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
