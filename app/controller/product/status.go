package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func (ctr *Controller) UpdateStatus(column dto.ProductStatusColumn, status bool) func(*gin.Context) {
	return func(ctx *gin.Context) {
		// Obtener productos
		if err := ctr.repositories.product.UpdateStatus(ctx.Param("id"), column, status); err != nil {
			server.InternalErrorException("Ocurrio un error", nil)
		}
		// Regresar data
		server.SuccessResponse(ctx)
	}
}
