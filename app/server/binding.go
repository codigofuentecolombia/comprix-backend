package server

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GinUriBinding[T any](ginContext *gin.Context, data T, message string) T {
	// Verificar si existe error al obtener la data
	if err := ginContext.ShouldBindUri(&data); err != nil {
		BadRequestException(message, data)
	}
	// Rgresar valor
	return data
}

func GinFormBinding[T any](ginContext *gin.Context, data T, message string) T {
	// Verificar si existe error al obtener la data
	if err := ginContext.ShouldBind(&data); err != nil {
		BadRequestException(message, data)
	}
	// Rgresar valor
	return data
}

func GinQueryAndFormBinding[T any](ginContext *gin.Context, data T) T {
	formErr := ginContext.ShouldBind(&data)
	queryErr := ginContext.ShouldBindQuery(&data)
	// Validar si al menos uno no tuvo error
	if formErr != nil && queryErr != nil {
		BadRequestException("Invalid form", data)
	}
	// Rgresar valor
	return data
}

func GinUriAndFormBinding[T any](ginContext *gin.Context, data T) T {
	uriErr := ginContext.ShouldBindUri(&data)
	formErr := ginContext.ShouldBind(&data)
	// Validar si al menos uno no tuvo error
	if formErr != nil && uriErr != nil {
		BadRequestException("Invalid form", data)
	}
	// Rgresar valor
	return data
}

func GinUriAndQueryBinding[T any](ginContext *gin.Context, data T) T {
	uriErr := ginContext.ShouldBindUri(&data)
	formErr := ginContext.ShouldBindQuery(&data)
	// Validar si al menos uno no tuvo error
	if formErr != nil && uriErr != nil {
		BadRequestException("Invalid form", data)
	}
	// Rgresar valor
	return data
}

func StructCopier[T any](dest T, source interface{}) T {
	// Verificar si existe error al obtener la data
	if err := copier.Copy(&dest, source); err != nil {
		InternalErrorException("Copy struct error", nil)
	}
	// Rgresar valor
	return dest
}
