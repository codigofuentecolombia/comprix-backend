package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctr *Controller) GetPaginated(ginContext *gin.Context) {
	order := "category_asc"
	bestPrice := true
	omitDisabled := true
	// Obtener respuesta
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repositories.pageProduct.GetPaginated(
			server.GinUriAndQueryBinding(ginContext, dto.ProductRepositoryFindParams{
				Order:        &order,
				Selects:      &dto.RepositoryGormSelections{Query: "page_products.*"},
				BestPrice:    &bestPrice,
				OmitDisabled: &omitDisabled,
				Preloads: &[]dto.RepositoryGormParams{
					{Query: "Page"},
					{Query: "Product"},
					{Query: "Product.Brand"},
					{Query: "Product.Category"},
				},
			}),
		),
	))
}
