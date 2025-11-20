package pages_jumbo

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
	var product dto.RetrievedProduct
	var hasStock bool
	var isNotFound bool
	var priceDetail dto.RetrievedProductPrice
	// Verificar el tipo de producto
	err := service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Esperar a que el contenido este cargado
		chromedp.ActionFunc(service.svc.MouseActions(".render-provider")),
		// Verificar si existe
		chromedp.Evaluate(`document.querySelector(".vtex-flex-layout-0-x-flexRow--row-opss-notfound") !== null`, &isNotFound),
	}, 20)
	// Terminar proceso si esta desactivado
	if isNotFound {
		return nil
	}
	// Obtener data
	err = service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Esperar a que el contenido este cargado
		chromedp.ActionFunc(service.svc.MouseActions(".vtex-store-components-3-x-container .vtex-flex-layout-0-x-flexCol--product-box .vtex-flex-layout-0-x-flexColChild--shelf-main-price-box div span div div")),
		// Esperar la carga
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.WaitUntilProductPriceIsLoaded()), &hasStock),
		// Esperar hasta que este listo
		chromedp.WaitVisible(`#succes-load`, chromedp.ByID),
		// Obtener datos de producto
		chromedp.Text(`.vtex-store-components-3-x-productNameContainer`, &product.Name, chromedp.ByQuery),
		chromedp.Text(`span.vtex-store-components-3-x-productBrandName`, &product.Brand, chromedp.ByQuery),
		chromedp.Text(`.view-conditions_more_info .view-conditions_descripcion`, &product.Description, chromedp.ByQuery),
		chromedp.Text(`span.vtex-product-identifier-0-x-product-identifier__value`, &product.Sku, chromedp.ByQuery),
		//
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetImagesScript()), &product.Images),
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetCategoriesScript()), &product.Categories),
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetPriceDetailScript()), &priceDetail),
	}, 40)
	//Verificar si existe un error
	if err != nil {
		return service_scrapper.HandleTotalTries(service.svc, dto.HandleScrapperTotalTriesParams[*dto.RetrievedProduct]{
			Err:             err,
			Url:             params.Url,
			Msg:             "No se pudo obtener el detalle de un producto",
			Callback:        service.GetProductDetail,
			CallbackArgs:    params,
			DefaultResponse: nil,
		})
	}
	// Setear valores
	// Actualizar daata
	product.Url = params.Url
	product.Sku = html.UnescapeString(product.Sku)
	product.Name = html.UnescapeString(product.Name)
	product.Brand = html.UnescapeString(product.Brand)
	product.PageID = service.svc.Page.ID
	product.HasStock = &hasStock
	product.OriginalPrice = priceDetail.Price
	product.OriginalDiscountPrice = priceDetail.DiscountPrice
	product.MinQuantityToApplyDiscount = priceDetail.MinQuantity
	// Verificar si tiene precio original
	if product.OriginalPrice == "" || product.OriginalPrice == "0" {
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
