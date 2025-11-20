package order_controller

import (
	"comprix/app/constants"
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ctr OrderController) GetDetail(ginContext *gin.Context) {
	id := ginContext.Param("id")
	var userID interface{} = nil
	authUser := ginContext.MustGet("authUser").(*dao.User)
	// Obteenr orden
	order, err := ctr.repositories.order.FindByID(id)
	// Validar si se pudo obtener
	if err != nil {
		server.InternalErrorException("No se pudo obtener el detalle de la orden", nil)
	}
	// Verificar los roles
	if authUser.Role.Name != constants.AdminRole {
		userID = authUser.ID
	}
	// Obtener producto
	items, err := ctr.repositories.order.FindAllByID(order.ID, userID)
	// Validar si se obtuvieron los items
	if err != nil {
		server.InternalErrorException("No se pudo obtener el detalle de la orden", nil)
	}
	// Obtener usuario
	user, err := ctr.repositories.user.FindByID(order.UserID)
	// Validar si se obtuvo el usuario
	if err != nil {
		server.InternalErrorException("No se pudo obtener el detalle de la orden", nil)
	}
	// Regresar
	server.Response(ginContext, http.StatusOK, dto.OrderDetail{
		User:            *user,
		Order:           *order,
		Items:           items,
		ShippingAddress: order.UserShippingAddress,
	})
}

func (ctr OrderController) GetStatistics(ginContext *gin.Context) {
	server.Response(ginContext, http.StatusOK, ctr.repositories.order.GetStatistics())
}
