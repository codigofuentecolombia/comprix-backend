package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"

	"github.com/gin-gonic/gin"
)

func PageRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/pages")
	controller := controller.InitPageController(config)
	// Listrar rutas
	router.GET("", controller.FindAll)
}
