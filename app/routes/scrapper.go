package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func ScrapperRoutes(config dto.Config, router *gin.RouterGroup) {
	// Crear controlador
	scrapperController := controller.NewScrapperController(&config)
	
	// Ruta protegida solo para administradores
	router.POST("/admin/scrapper/start", server.AdminMiddleware(config), scrapperController.StartScrapper)
}
