package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"

	"github.com/gin-gonic/gin"
)

func HealthRoutes(config dto.Config, router *gin.RouterGroup) {
	healthController := controller.NewHealthController(&config)
	
	// Endpoint público para verificar la conexión a la BD
	router.GET("/health/db", healthController.CheckDB)
}
