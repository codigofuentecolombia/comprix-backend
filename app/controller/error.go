package controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorController struct {
	config     *dto.Config
	repository *repositories.ErrorRepository
}

func InitErrorController(config *dto.Config) ErrorController {
	return ErrorController{
		config:     config,
		repository: repositories.InitErrorRepository(config.GormDB),
	}
}

func (ctr ErrorController) GetPaginated(ginContext *gin.Context) {
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repository.GetPaginated(
			server.GinQueryAndFormBinding(ginContext, dto.ErrorRepositoryFindParams{
				Pagination: server.GinQueryAndFormBinding(ginContext, dto.Pagination[dao.Error]{}),
				// Cargar relaciones
				Preloads: &[]dto.RepositoryGormParams{
					{Query: "Page"},
				},
			}),
		),
	))
}
