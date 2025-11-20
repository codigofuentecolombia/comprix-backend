package controller

import (
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	repository_page_product "comprix/app/repositories/page-product"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageProductController struct {
	config          *dto.Config
	repository      *repositories.ProductRepository
	pageProductRepo *repository_page_product.Repository
}

func InitPageProductController(config *dto.Config) PageProductController {
	return PageProductController{
		config:          config,
		repository:      repositories.InitProductRepository(config.GormDB),
		pageProductRepo: repository_page_product.InitRepository(config.GormDB),
	}
}

func (ctr *PageProductController) Update(ginContext *gin.Context) {
	err := ctr.repository.Update(
		server.GinUriAndFormBinding(ginContext, dto.UpdatePageProduct{}),
	)
	// Verificar si ocurrio un error
	if err != nil {
		server.InternalErrorException("No se pudo actualizar el producto", nil)
	}
	// Regresar respuesta exitosa
	server.SuccessResponse(ginContext)
}

func (ctr *PageProductController) GetDetail(ginContext *gin.Context) {
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repository.Find(
			server.GinUriAndFormBinding(ginContext, dto.ProductRepositoryFindParams{
				Preloads: &[]dto.RepositoryGormParams{
					{Query: "Product"},
				},
			}),
		),
	))
}

func (ctr *PageProductController) GetCheapest(ginContext *gin.Context) {
	bestPrice := true
	// Obtener formulario
	form := server.GinUriAndFormBinding(ginContext, dto.ProductRepositoryFindParams{
		Selects: &dto.RepositoryGormSelections{Query: "page_products.*"},
		Preloads: &[]dto.RepositoryGormParams{
			{Query: "Page"},
			{Query: "Product"},
			{Query: "Product.Brand"},
			{Query: "Product.Category"},
		},
	})
	// Obtener producto
	product := server.EntityNotFound(ctr.pageProductRepo.Find(form))
	// Verificar si existe id
	if form.ID != nil {
		// Vaciar formulario
		form.ID = nil
		form.ExcludeID = &product.ID
		form.ProductID = &product.MainProductID
		form.BestPagePrice = &bestPrice
		// Mostrar solo pagina
		form.Preloads = &[]dto.RepositoryGormParams{
			{Query: "Page"},
		}
		// Obtener
		products := server.EntityNotFound(ctr.pageProductRepo.FindAll(form))
		product.OtherStores = &products
	}
	// Obtener respuesta
	server.Response(ginContext, http.StatusOK, product)
}
