package controller

import (
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	config *dto.Config
}

func NewHealthController(config *dto.Config) HealthController {
	return HealthController{
		config: config,
	}
}

func (ctr HealthController) CheckDB(ginContext *gin.Context) {
	// Intentar contar usuarios
	userRepo := repositories.InitUserRepository(ctr.config.GormDB)
	
	// Ejecutar una query simple
	var count int64
	err := ctr.config.GormDB.Table("users").Count(&count).Error
	
	if err != nil {
		server.Response(ginContext, 500, gin.H{
			"status": "error",
			"message": "No se pudo conectar a la base de datos",
			"error": err.Error(),
			"dsn_configured": ctr.config.Settings.Database.Dsn != "",
		})
		return
	}
	
	server.Response(ginContext, 200, gin.H{
		"status": "ok",
		"message": "Conexión a BD exitosa",
		"users_count": count,
		"database_connected": true,
	})
	
	_ = userRepo
}
