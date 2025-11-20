package auth_controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func (ctr AuthController) Register(ginContext *gin.Context) {
	// Obtener parametros validados
	form := server.GinFormBinding(ginContext, dto.AuthRegister{}, "Invalid form")
	// Obtener usuario
	_, err := ctr.repository.FindByUsernameNoRelation(form.Username)
	// Verificar si ocurrio error
	if err == nil {
		server.BadRequestException("Usuario en uso", nil)
	}
	// Crear usuario
	_, err = ctr.service.HandleLogin(dao.User{
		Email:       form.Username,
		RoleID:      2,
		LastName:    form.LastName,
		Username:    form.Username,
		Password:    form.Password,
		FirstName:   form.FirstName,
		PhoneNumber: form.PhoneNumber,
	})
	// Validar si ocurrio un error
	if err != nil {
		server.InternalErrorException("No se pudo registrar al usuario", nil)
	}
	// Regresa respuesta
	server.SuccessResponse(ginContext)
}
