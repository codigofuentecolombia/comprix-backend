package auth_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctr AuthController) Socialite(sType dto.SocialiteType) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		server.Response(ginContext, http.StatusOK, dto.AuthSocialiteResponse{
			Url: ctr.service.GenerateAuthConfigByType(sType).AuthCodeURL("state"),
		})
	}
}

func (ctr AuthController) SocialiteCallback(sType dto.SocialiteType) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Generar configuracion de autentificacion
		authConfig := ctr.service.GenerateAuthConfigByType(sType)
		// Verificar si se encontro el codigo
		request := server.GinFormBinding(ginContext, dto.AuthRequestCode{}, "Código no recibido")
		// Intercambiar el código por un token de acceso
		token, err := authConfig.Exchange(context.Background(), request.Code)
		// Verficar si se pudo descifrar token
		if err != nil {
			server.InternalErrorException("No se pudo intercambiar el código por un token", nil)
		}
		// Obtener información del usuario desde Facebook
		client := authConfig.Client(context.Background(), token)
		// Consultar detalles de usuario
		resp, err := client.Get(ctr.service.GetDetailUriByType(sType))
		// Verificar si ocurrio un error
		if err != nil {
			server.InternalErrorException("No se pudo obtener la información del usuario", nil)
		}
		// Cerrar respuesta
		defer resp.Body.Close()
		// Obtener usuario
		user, err := ctr.service.DecodeSocialiteUser(sType, resp.Body)
		// Validar si  ocurrio un error
		if err != nil {
			server.InternalErrorException("No se pudo obtener la información del usuario", nil)
		}
		//
		authRes, err := ctr.service.HandleLogin(*user)
		// Validar si ocurrio un error
		if err != nil {
			server.InternalErrorException("No se pudo obtener la información del usuario", nil)
		}
		// Regresar respuesta
		server.Response(ginContext, http.StatusOK, authRes)
	}
}
