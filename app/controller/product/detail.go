package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctr *Controller) GetCheapest(omitDisabled bool) func(*gin.Context) {
	return func(ginContext *gin.Context) {
		// Obtener producto de pagina
		currentPageProduct := server.EntityNotFound(ctr.repositories.pageProduct.Find(
			server.GinUriAndFormBinding(ginContext, dto.ProductRepositoryFindParams{
				Selects: &dto.RepositoryGormSelections{Query: "page_products.main_product_id"},
			}),
		))
		// Obtener los mejores precios por pagina
		bestPrice := true
		// Obtener
		products := server.EntityNotFound(ctr.repositories.pageProduct.FindAll(dto.ProductRepositoryFindParams{
			Selects: &dto.RepositoryGormSelections{Query: []string{
				"page_products.id",
				"page_products.price",
				"page_products.images",
				"page_products.page_id",
				"page_products.discount_price",
				"page_products.min_quantity_to_apply_discount",
			}},
			Preloads:      &[]dto.RepositoryGormParams{{Query: "Page"}},
			ProductID:     &currentPageProduct.MainProductID,
			OmitDisabled:  &omitDisabled,
			BestPagePrice: &bestPrice,
		}))
		// Verificar si existen productos
		if len(products) == 0 {
			server.InternalErrorException("Could not found entity", nil)
		}
		// Crear producto
		tmpProduct := products[0]
		// Iterar productos omitiendo index 0
		for i := 1; i < len(products); i++ {
			product := products[i]
			isCheaper := false
			// Verificar si tiene menor precio
			if product.DiscountPrice > 0 {
				if tmpProduct.DiscountPrice > 0 {
					isCheaper = product.DiscountPrice < tmpProduct.DiscountPrice
				} else {
					isCheaper = product.DiscountPrice < tmpProduct.Price
				}
			} else {
				if tmpProduct.DiscountPrice > 0 {
					isCheaper = product.Price < tmpProduct.DiscountPrice
				} else {
					isCheaper = product.Price < tmpProduct.Price
				}
			}
			// Verificar si es mas barato para sustituir producto
			if isCheaper {
				tmpProduct = product
			}
		}
		// Obtener el detalle extenso del producto mas barato
		product := server.EntityNotFound(ctr.repositories.pageProduct.Find(dto.ProductRepositoryFindParams{
			ID:      &tmpProduct.ID,
			Selects: &dto.RepositoryGormSelections{Query: "page_products.*"},
			Preloads: &[]dto.RepositoryGormParams{
				{Query: "Page"},
				{Query: "Product"},
				{Query: "Product.Brand"},
				{Query: "Product.Category"},
			},
		}))
		// Agregar productos
		product.OtherStores = &products
		// Obtener respuesta
		server.Response(ginContext, http.StatusOK, product)
	}
}
