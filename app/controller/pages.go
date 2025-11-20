package controller

import (
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageController struct {
	config     dto.Config
	repository repositories.PageRepository
}

func InitPageController(config dto.Config) PageController {
	return PageController{
		config:     config,
		repository: repositories.InitPageRepository(config.GormDB),
	}
}

func (ctr PageController) FindAll(ginContext *gin.Context) {
	entities, err := ctr.repository.FindAll()
	// Validar si hubo error
	if err != nil {
		server.InternalErrorException("Ocurrio un error", err)
	}
	// Regresar data
	server.Response(ginContext, http.StatusOK, entities)
}
