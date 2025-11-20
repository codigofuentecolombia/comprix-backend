package user_controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	service_auth "comprix/app/services/auth"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	config     dto.Config
	service    service_auth.Service
	repository repositories.UserRepository
}

func InitController(config dto.Config) Controller {
	return Controller{
		config:     config,
		service:    service_auth.InitService(config),
		repository: repositories.InitUserRepository(config.GormDB),
	}
}

func (c Controller) GetPaginated(ginContext *gin.Context) {
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
	params := dto.GetUsersParams{
		Pagination: dto.Pagination[dao.User]{
			Index: page,
			Limit: limit,
		},
	}
	// Validar paginacion
	params.Pagination.Validate()
	// Obtener productos
	products, err := c.repository.GetPaginated(params)
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, products)
}
