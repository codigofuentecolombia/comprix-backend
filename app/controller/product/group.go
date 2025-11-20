package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func (ctr *Controller) GroupExisting(ginContext *gin.Context) {
	err := ctr.repositories.pageProduct.GroupExisting(
		server.GinUriBinding(ginContext, dto.GroupPageProduct{}, "No se pudo obtener los parametros de la petición"),
	)
	// Verificar si ocurrio error
	if err != nil {
		server.InternalErrorException("No se pudo agrupar los productos especificados", nil)
	}
	// Regresar exitosa respuesta
	server.SuccessResponse(ginContext)
}
