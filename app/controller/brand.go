package controller

import (
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BrandController struct {
	config     *dto.Config
	repository *repositories.BrandhRepository
}

func InitBrandController(config *dto.Config) BrandController {
	return BrandController{
		config:     config,
		repository: repositories.InitBrandhRepository(config.GormDB),
	}
}

func (ctr BrandController) GetAll(ginContext *gin.Context) {
	withProducts := true
	//
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repository.FindAll(
			server.GinFormBinding(ginContext, dto.BrandRepositoryFindParams{HasProducts: &withProducts}, "No se pudo obtener los parametros"),
		),
	))
}
