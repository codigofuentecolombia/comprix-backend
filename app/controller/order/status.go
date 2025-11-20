package order_controller

import (
	"comprix/app/domain/dao"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func (ctr OrderController) UpdateStatus(status dao.OrderStatus) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		id := ginContext.Param("id")
		err := ctr.repositories.order.UpdateStatus(id, status)
		// Validar si se pudo obtener
		if err != nil {
			server.InternalErrorException("No se pudo actualizar el status de la orden", nil)
		}
		// Regresar
		server.SuccessResponse(ginContext)
	}
}
