package controller

import (
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	"comprix/app/server"
	"comprix/app/services/auth"
	"github.com/gin-gonic/gin"
)

type DebugPasswordController struct {
	config     *dto.Config
	repository *repositories.UserRepository
}

func NewDebugPasswordController(config *dto.Config) *DebugPasswordController {
	return &DebugPasswordController{
		config:     config,
		repository: repositories.NewUserRepository(config.Db),
	}
}

func (ctr *DebugPasswordController) TestPassword(c *gin.Context) {
	var form struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&form); err != nil {
		server.BadRequestException("Datos inválidos", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Buscar usuario
	user, err := ctr.repository.FindByUsernameNoRelation(form.Email)
	if err != nil {
		server.BadRequestException("Usuario no encontrado", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Verificar contraseña
	isValid := auth.CheckPasswordHash(form.Password, user.Password)

	// Retornar información de debug
	server.OkResponse(map[string]interface{}{
		"user_id":          user.Id,
		"email":            user.Email,
		"password_hash":    user.Password,
		"password_tested":  form.Password,
		"password_matches": isValid,
	})
}
