package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func ErrorRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/errors", server.ValidateJwt(config), server.Role("admin"))
	controller := controller.InitErrorController(&config)
	// Listrar rutas
	router.GET("", controller.GetPaginated)
}
