package pages_masonline

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
	// Obtener data
	err := service.svc.InitAndRunActions(params.Url, []chromedp.Action{
		// Obtener datos de producto
		chromedp.ActionFunc(service.svc.MouseActions(".vtex-store-components-3-x-productBrand")),
		// Obtener valores
		chromedp.Text(`.vtex-store-components-3-x-productBrand`, &product.Name, chromedp.ByQuery),
		chromedp.Text(`.vtex-flex-layout-0-x-flexRow--product-view-pricesandbutton-containerchild .valtech-gdn-dynamic-product-1-x-weighablePriceWrapper .valtech-gdn-dynamic-product-1-x-weighableSavings`, &product.OriginalPrice, chromedp.ByQuery),
		chromedp.ActionFunc(service.svc.GetSecuredValue(`.vtex-store-components-3-x-productDescriptionText`, &product.Description)),
		// Atributos
		chromedp.AttributeValue(`meta[property="product:sku"]`, "content", &product.Sku, nil),
		chromedp.AttributeValue(`meta[property="product:brand"]`, "content", &product.Brand, nil),
		chromedp.AttributeValue(`meta[property="product:price:amount"]`, "content", &product.OriginalDiscountPrice, nil),
		// Obtener product
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetImagesScript()), &product.Images),
		chromedp.Evaluate(fmt.Sprintf("(() => {%s})()", service.GetCategoriesScript()), &product.Categories),
	}, 200)
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
	// Actualizar daata
	product.Url = params.Url
	product.Sku = html.UnescapeString(product.Sku)
	product.Name = html.UnescapeString(product.Name)
	product.Brand = html.UnescapeString(product.Brand)
	product.PageID = service.svc.Page.ID
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
