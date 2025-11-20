package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func CalendarRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/calendar")
	controller := controller.InitCalendarController(config)
	// Listrar rutas
	router.GET("", controller.GetAll)
	router.PUT("", server.ValidateJwt(config), server.Role("admin"), controller.Update)
	router.POST("", server.ValidateJwt(config), server.Role("admin"), controller.Create)
	router.DELETE(":id", server.ValidateJwt(config), server.Role("admin"), controller.Delete)
}
