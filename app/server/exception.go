package server

import (
	"comprix/app/domain/dto"
	"comprix/app/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Exception struct {
	msg        string
	payload    any
	statusCode int
}

// Funcion para emitir excepcion con status code 400
func BadRequestException(msg string, data any) {
	emitException(http.StatusBadRequest, msg, data)
}

// Funcion para emitir excepcion con status code 401
func UnauthorizedException(msg string, data any) {
	emitException(http.StatusUnauthorized, msg, data)
}

// Funcion para emitir excepcion con status code 403
func ForbiddenException(msg string, data any) {
	emitException(http.StatusForbidden, msg, data)
}

// Funcion para emitir excepcion con status code 404
func NotFoundException(msg string, data any) {
	emitException(http.StatusNotFound, msg, data)
}

// Funcion para emitir excepcion con status code 500
func InternalErrorException(msg string, data any) {
	emitException(http.StatusInternalServerError, msg, data)
}

func NotAcceptableException(msg string, data any) {
	emitException(http.StatusNotAcceptable, msg, data)
}

// Emitir exception
func emitException(statusCode int, msg string, payload any) {
	panic(&Exception{
		msg:        msg,
		payload:    payload,
		statusCode: statusCode,
	})
}

// Manejar excepciones de gin
func HandleException(config dto.Config) gin.HandlerFunc {
	// Crear loger de error
	serverLog := logger.Create(config.Settings.Paths.Logs, "server/errors")
	// Regresar handler
	return func(ginContext *gin.Context) {
		// Manejar el panico
		defer func() {
			if err := recover(); err != nil {
				// Crear log
				log := serverLog.WithFields(logrus.Fields{
					// "trace":   string(debug.Stack()),
					// "request": GetRequestDetail(ginContext),
				})
				//Validar si el tipo de error es HttpException
				if e, ok := err.(*Exception); ok {
					response := map[string]any{
						"msg": e.msg,
					}
					// Verificar si existe payload
					if e.payload != nil {
						response["payload"] = e.payload
					}
					// Regresar respuesta
					Response(ginContext, e.statusCode, response)
					// Mostrar log
					log.Error(e)
				} else {
					// Obtener error
					e := err.(error)
					// Responder
					Response(ginContext, http.StatusInternalServerError, map[string]any{
						"msg": e.Error(),
					})
					// Mostrar log
					log.Error(e)
				}
				// Terminar
				ginContext.Abort()
			}
		}()
		// Llama a los siguientes middlewares y handlers
		ginContext.Next()
	}

}
