package controller

import (
	"comprix/app/domain/dto"
	repository_category "comprix/app/repositories/category"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	config     *dto.Config
	repository *repository_category.Repository
}

func InitCategoryController(config *dto.Config) CategoryController {
	return CategoryController{
		config:     config,
		repository: repository_category.InitRepository(config.GormDB),
	}
}

func (ctr CategoryController) GetAll(ginContext *gin.Context) {
	onlyParents := true
	//
	server.Response(ginContext, http.StatusOK, server.EntityNotFound(
		ctr.repository.FindAll(
			server.GinUriAndQueryBinding(ginContext, dto.CategoryRepositoryFindParams{
				OnlyParents: &onlyParents,
				Preloads: &[]dto.RepositoryGormParams{
					{Query: "Subcategories"},
					{Query: "Subcategories.Subcategories"},
				},
			}),
		),
	))
}
