package order_controller

import (
	"comprix/app/constants"
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/server"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (c OrderController) GetPaginated(ginContext *gin.Context) {
	authUser := ginContext.MustGet("authUser").(*dao.User)
	// Valores predeterminados
	defaultPage := 1
	defaultLimit := 10
	// Obtener los valores de la URL
	pageStr := ginContext.Query("page")   // Obtiene el valor de "pagina"
	limitStr := ginContext.Query("limit") // Obtiene el valor de "limite"
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
	params := dto.GetOrdersParams{
		Pagination: dto.Pagination[dao.Order]{
			Index: page,
			Limit: limit,
		},
	}
	// Verificar los roles
	if authUser.Role.Name != constants.AdminRole {
		params.UserID = authUser.ID
	}
	// Validar paginacion
	params.Pagination.Validate()
	// Obtener productos
	products, err := c.repositories.order.GetPaginated(params)
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}
