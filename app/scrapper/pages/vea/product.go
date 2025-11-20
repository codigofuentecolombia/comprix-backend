package pages_vea

import (
	"comprix/app/domain/dto"
	service_scrapper "comprix/app/services/scrapper"
	"comprix/app/utils"
	"fmt"
	"html"

	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

func (service Service) GetProductDetail(params dto.ScrapperParams) *dto.RetrievedProduct {
	// Mostrar log
	service.svc.CreateDebugLog("Iniciando scraping del producto.", logrus.Fields{"url": params.Url})
	// Definir producto
	var dummy any
	var product dto.RetrievedProduct
	// Obtener data
	err := service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Esperar a que el contenido este cargado
		chromedp.ActionFunc(service.svc.MouseActions(".vtex-flex-layout-0-x-flexCol--product-box")),
		// Esperar la carga
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.WaitUntilProductPriceIsLoaded()), &dummy),
		// Esperar hasta que este listo
		chromedp.WaitVisible(`#succes-load`, chromedp.ByID),
		// Obtener datos de producto
		chromedp.Text(`#priceContainer`, &product.OriginalDiscountPrice, chromedp.NodeVisible),
		chromedp.Text(`.view-conditions_more_info .view-conditions_descripcion`, &product.Description, chromedp.ByQuery),
		chromedp.Text(`.vtex-flex-layout-0-x-flexRow--product-main .vtex-store-components-3-x-productBrand`, &product.Name, chromedp.ByQuery),
		chromedp.Text(`.vtex-flex-layout-0-x-flexRow--product-main .vtex-store-components-3-x-productBrandName`, &product.Brand, chromedp.ByQuery),
		chromedp.Text(`.vtex-flex-layout-0-x-flexRow--product-main .vtex-product-identifier-0-x-product-identifier__value`, &product.Sku, chromedp.ByQuery),
		// Obtener product
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetImagesScript()), &product.Images),
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetCategoriesScript()), &product.Categories),
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetOriginalPriceScript()), &product.OriginalPrice),
	}, 120)
	//Verificar si existe un error
	if err != nil {
		return service_scrapper.HandleTotalTries(service.svc, dto.HandleScrapperTotalTriesParams[*dto.RetrievedProduct]{
			Err:             err,
			Url:             params.Url,
			Msg:             "Ocurrio un error al tratar de obtener el detalle de un producto",
			Callback:        service.GetProductDetail,
			CallbackArgs:    params,
			DefaultResponse: nil,
		})
	}
	// Actualizar daata
	product.Url = params.Url
	product.Sku = html.UnescapeString(product.Sku)
	product.Name = html.UnescapeString(product.Name)
	product.Brand = html.UnescapeString(product.Brand)
	product.PageID = service.svc.Page.ID
	// Verificar si tiene precio original
	if product.OriginalPrice == "0" || product.OriginalPrice == "" {
		product.OriginalPrice = product.OriginalDiscountPrice
		product.OriginalDiscountPrice = ""
	}
	// Obtener precio
	price, err := utils.CleanCurrencyFormat(product.OriginalPrice)
	discountPrice, _ := utils.CleanCurrencyFormat(product.OriginalDiscountPrice)
	// Verificar si tiene precio
	if err != nil {
		return service_scrapper.HandleTotalTries(service.svc, dto.HandleScrapperTotalTriesParams[*dto.RetrievedProduct]{
			Err:             err,
			Url:             params.Url,
			Msg:             fmt.Sprintf("No fue posible limpiar el formato del precio del producto: %s", product.OriginalPrice),
			Callback:        service.GetProductDetail,
			CallbackArgs:    params,
			DefaultResponse: nil,
		})
	}
	// Actualizar data
	product.Price = price
	product.DiscountPrice = discountPrice
	// Imprimir
	return &product
}
