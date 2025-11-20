package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"

	"github.com/gin-gonic/gin"
)

func BrandsRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/brands")
	controller := controller.InitBrandController(&config)
	// Listrar rutas
	router.GET("", controller.GetAll)
}
