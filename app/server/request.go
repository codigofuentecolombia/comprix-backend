package server

import (
	"comprix/app/domain/dto"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRequestDetail(ginContext *gin.Context) dto.ServerRequest {
	// Obtener headers del request
	headers := []string{}
	// Iterar los headers del request
	for key, value := range ginContext.Request.Header {
		headers = append(headers, key+": "+strings.Join(value, " "))
	}
	// Regresar el detalle
	return dto.ServerRequest{
		URL:     ginContext.Request.URL.Path,
		Method:  ginContext.Request.Method,
		Params:  ginContext.Request.URL.RawQuery,
		Headers: headers,
	}
}
