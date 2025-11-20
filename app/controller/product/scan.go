package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ctr *Controller) Scan(ginContext *gin.Context) {
	// Obtener datos enviados
	params := server.GinUriAndFormBinding(ginContext, dto.ScanPageProduct{})
	// Declarar variables
	var page dto.IScrapperService
	// Obtener producto dependiendo de la pagina
	if strings.Contains(params.Url, "jumbo.com.ar") {
		page = ctr.pages.jumbo
	} else if strings.Contains(params.Url, "hiperlibertad.com.ar") {
		page = ctr.pages.hiperlibertad
	} else if strings.Contains(params.Url, "vea.com.ar") {
		page = ctr.pages.vea
	} else if strings.Contains(params.Url, "masonline.com.ar") {
		page = ctr.pages.masonline
	} else if strings.Contains(params.Url, "carrefour.com.ar") {
		page = ctr.pages.carrefour
	} else {
		server.BadRequestException("La url no pertenece a una pagina permetida", nil)
	}
	// Obtener categoria
	categoryResponse := ctr.repositories.category.Find(dto.CategoryRepositoryFindParams{
		ID:       &params.Product,
		ParentID: &params.Subcategory,
	})
	// Verificar si se encontro
	if categoryResponse.Error != nil {
		server.InternalErrorException("Error al buscar la categoria", nil)
	}
	// Obtener producto
	retrievedProduct := page.GetProductDetail(dto.ScrapperParams{
		Url:        params.Url,
		Categories: []string{categoryResponse.Data.Name},
	})
	// Agregar categorias al producto
	retrievedProduct.Categories = []string{categoryResponse.Data.Name}
	// Crear producto
	page.CreateOrUpdateProduct(retrievedProduct)
	// Obtener error
	shouldOrder := true
	errResponse := ctr.repositories.err.Find(dto.ErrorRepositoryFindParams{
		Url:             &params.Url,
		ShouldOrderDesc: &shouldOrder,
	})
	// Verificar si no ocurrio error
	if errResponse.Error == nil {
		server.InternalErrorException("Error al escanear producto", errResponse.Data.Message)
	}
	// Finalizar peticion
	server.SuccessJsonResponse(ginContext)
}
