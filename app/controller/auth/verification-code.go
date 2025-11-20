package auth_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"
	service_email "comprix/app/services/email"

	"github.com/gin-gonic/gin"
)

func (ctr AuthController) SendVerificationCode(ginContext *gin.Context) {
	// Obtener parametros validados
	form := server.GinFormBinding(ginContext, dto.AuthSendVerificationCode{}, "Invalid form")
	// Obtener usuario
	user, err := ctr.repository.FindByEmail(form.Email)
	// Verificar si ocurrio error
	if err != nil {
		server.BadRequestException("Credenciales incorrecta", nil)
	}
	// Crear codigo de verificacion
	vCode, err := ctr.service.GenerateVerificationCode(*user)
	// Verificar si hubo error
	if err != nil {
		server.BadRequestException("Credenciales incorrecta", nil)
	}
	// Enviar
	mailer := service_email.InitService(ctr.config)
	// Enviar
	mailer.VerificationCode(*user, vCode)
	// Regresar respuesta
	server.SuccessResponse(ginContext)
}

func (ctr AuthController) VerificationCode(ginContext *gin.Context) {
	// Verificar
	if ctr.service.VerifyCode(server.GinUriAndFormBinding(ginContext, dto.AuthVerificationCode{})) {
		server.BadRequestException("Credenciales incorrecta", nil)
	}
	// Regresar
	server.SuccessResponse(ginContext)
}
