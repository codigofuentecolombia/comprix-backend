package auth_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func (ctr AuthController) Login(ginContext *gin.Context) {
	// Obtener parametros validados
	form := server.GinFormBinding(ginContext, dto.AuthLogin{}, "Invalid form")
	// Obtener usuario
	user, err := ctr.repository.FindByUsername(form.Username)
	// Verificar si ocurrio error
	if err != nil {
		server.BadRequestException("Credenciales incorrecta", nil)
	}
	// Validar si las contraseñas coinciden
	if ctr.service.CheckPasswordHash(form.Password, user.Password) != nil {
		server.BadRequestException("Credenciales incorrecta", nil)
	}
	// Obtener jwt
	token, err := ctr.service.GenerateToken(*user)
	// Verificar si ocurrio un error al generar el jwt
	if err != nil {
		server.InternalErrorException("No se pudo generar el JWT", nil)
	}
	// Verificar si tiene el correo verificado
	if !user.IsVerified {
		server.NotAcceptableException("No se ha verificado el correo", nil)
	}
	// Regresar respuesta
	server.Response(ginContext, 200, dto.AuthResponse{
		User:  *user,
		Token: token,
	})
}
