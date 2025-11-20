package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctr *Controller) GetAll(ginContext *gin.Context) {
	order := "admin"
	omitDisabled := true
	// Obtener respuesta
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repositories.pageProduct.GetPaginated(
			server.GinUriAndQueryBinding(ginContext, dto.ProductRepositoryFindParams{
				Order:        &order,
				Selects:      &dto.RepositoryGormSelections{Query: "page_products.*"},
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

func (ctr *Controller) GetRecommended(ginContext *gin.Context) {
	bestPrice := true
	onlyRecommended := true
	// Obtener respuesta
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repositories.pageProduct.FindAll(
			server.GinUriAndQueryBinding(ginContext, dto.ProductRepositoryFindParams{
				Selects:         &dto.RepositoryGormSelections{Query: "page_products.*"},
				BestPrice:       &bestPrice,
				OnlyRecommended: &onlyRecommended,
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

func (ctr *Controller) GetWithDiscount(ginContext *gin.Context) {
	bestPrice := true
	onlyWithDiscount := true
	// Obtener respuesta
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repositories.pageProduct.FindAll(
			server.GinUriAndQueryBinding(ginContext, dto.ProductRepositoryFindParams{
				Selects:          &dto.RepositoryGormSelections{Query: "page_products.*"},
				BestPrice:        &bestPrice,
				OnlyWithDiscount: &onlyWithDiscount,
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

func (ctr *Controller) GetDisabled(ginContext *gin.Context) {
	flag := true
	bestPrice := true
	// Obtener respuesta
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repositories.pageProduct.FindAll(
			server.GinUriAndQueryBinding(ginContext, dto.ProductRepositoryFindParams{
				Selects:      &dto.RepositoryGormSelections{Query: "page_products.*"},
				BestPrice:    &bestPrice,
				OnlyDisabled: &flag,
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
