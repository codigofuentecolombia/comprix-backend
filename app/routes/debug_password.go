package routes

import (
	"comprix/app/controller"
	"comprix/app/domain/dto"
	"github.com/gin-gonic/gin"
)

func DebugPasswordRoutes(config dto.Config, router *gin.RouterGroup) {
	ctrl := controller.NewDebugPasswordController(&config)

	// Endpoint público temporal para debug
	router.POST("/debug/test-password", ctrl.TestPassword)
}
