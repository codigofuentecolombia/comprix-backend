package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"

	"github.com/gin-gonic/gin"
)

func CategoryRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/categories")
	controller := controller.InitCategoryController(&config)
	// Listrar rutas
	router.GET("", controller.GetAll)
}
