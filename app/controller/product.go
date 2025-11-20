package controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	config dto.Config
}

func InitProductController(config dto.Config) Controller {
	return Controller{
		config: config,
	}
}

func (c Controller) GetByID(ginContext *gin.Context) {
	repository := repositories.InitPageProductRepository(c.config.GormDB)
	// Obtener productos
	products, err := repository.GetByID(ginContext.Param("id"))
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}

func (c Controller) GetOutstanding(ginContext *gin.Context) {
	repository := repositories.InitPageProductRepository(c.config.GormDB)
	// Obtener productos
	products, err := repository.GetAllWithLimit(8)
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}

func (c Controller) GetAll(ginContext *gin.Context) {
	repository := repositories.InitPageProductRepository(c.config.GormDB)
	// Obtener productos
	products, err := repository.GetAll(dto.GetProductsParams{})
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}

func (c Controller) GetPaginated(ginContext *gin.Context) {
	// Valores predeterminados
	defaultPage := 1
	defaultLimit := 10
	// Obtener los valores de la URL
	search := ginContext.Query("search")
	pageStr := ginContext.Query("page") // Obtiene el valor de "pagina"
	typeStr := ginContext.Query("type")
	orderBy := ginContext.Query("order_by")
	category := ginContext.Query("category")
	limitStr := ginContext.Query("limit") // Obtiene el valor de "limite"
	brandsIds := ginContext.QueryArray("branch_ids")
	// Convertir los valores a enteros con valores predeterminados si están vacíos
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = defaultLimit
	}
	// Params
	params := dto.GetProductsParams{
		Type:       dto.ProductsType(typeStr),
		OrderBy:    &orderBy,
		BranchIds:  &brandsIds,
		CategoryID: category,
		Pagination: dto.Pagination[dao.PageProduct]{
			Index:  page,
			Limit:  limit,
			Search: search,
		},
	}
	// Validar paginacion
	params.Pagination.Validate()
	//
	repository := repositories.InitPageProductRepository(c.config.GormDB)
	// Obtener productos
	products, err := repository.GetPaginated(params)
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}

func (c Controller) SetAsRecommended(ginContext *gin.Context) {
	// Obtener productos
	if err := repositories.SetAsRecommended(c.config.GormDB, ginContext.Param("id")); err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}

func (c Controller) UnsetAsRecommended(ginContext *gin.Context) {
	// Obtener productos
	if err := repositories.UnsetAsRecommended(c.config.GormDB, ginContext.Param("id")); err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}

func (c Controller) GetRecommended(ginContext *gin.Context) {
	repository := repositories.InitPageProductRepository(c.config.GormDB)
	// Obtener productos
	products, err := repository.GetRecommended()
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}

func (c Controller) SetAsInDiscount(ginContext *gin.Context) {
	// Obtener productos
	if err := repositories.SetAsInDiscount(c.config.GormDB, ginContext.Param("id")); err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}

func (c Controller) UnsetAsInDiscount(ginContext *gin.Context) {
	// Obtener productos
	if err := repositories.UnsetAsInDiscount(c.config.GormDB, ginContext.Param("id")); err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.SuccessResponse(ginContext)
}

func (c Controller) GetWithDiscount(ginContext *gin.Context) {
	repository := repositories.InitPageProductRepository(c.config.GormDB)
	// Obtener productos
	products, err := repository.GetWithDiscount()
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}
