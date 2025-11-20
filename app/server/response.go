package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Responder con data
func Response(ginContext *gin.Context, code int, data any) {
	ginContext.JSON(code, data)
}

// Responder con mensaje
func MessageResponse(ginContext *gin.Context, code int, msg string) {
	Response(ginContext, code, map[string]string{"msg": msg})
}

// Responder con exito
func SuccessResponse(ginContext *gin.Context) {
	Response(ginContext, http.StatusOK, map[string]bool{"success": true})
}

// Responder con exito
func SuccessJsonResponse(ginContext *gin.Context) {
	Response(ginContext, http.StatusOK, map[string]bool{"success": true})
}

// Responder la creacion
func CreateEntityResponse(ginContext *gin.Context, entity any) {
	Response(ginContext, http.StatusCreated, entity)
}
