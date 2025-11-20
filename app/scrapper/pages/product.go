package pages

import "comprix/app/domain/dto"

func (scrapper *ScrapPage) CreateOrUpdateProduct(link string, categories []string) {
	// Obtener los detalles del producto
	pageProduct := scrapper.svc.GetProductDetail(dto.ScrapperParams{Url: link, Categories: categories})
	// Verificar si existe producto
	if pageProduct != nil {
		// Verificar si hay categorias
		if len(categories) > 0 {
			pageProduct.Categories = categories
		}
		// Verificar si no se desactivara los productos
		if !scrapper.disableProduct {
			pageProduct.HasStock = nil
		}
	}
	// Crear o actualizar producto
	scrapper.svc.CreateOrUpdateProduct(pageProduct)
}

func RetrieveSingleProduct(conf *dto.Config, page dto.InitScrapPageFn, url string, categories []string) {
	scrapper, err := InitService(conf, page, 1)
	// Verificar si hubo error para finalizar funcion
	if err != nil {
		return
	}
	// Obtener producto
	scrapper.CreateOrUpdateProduct(url, categories)
}
