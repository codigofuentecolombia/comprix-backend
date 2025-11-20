package routes

import (
	order_controller "comprix/app/controller/order"
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/orders")
	controller := order_controller.InitOrderController(config)
	// Crear orden
	router.POST("", server.ValidateJwt(config), controller.New)
	router.GET("", server.ValidateJwt(config), controller.GetAll)
	router.GET(":id", server.ValidateJwt(config), controller.GetDetail)
	router.GET("paginated", server.ValidateJwt(config), controller.GetPaginated)
	router.GET("statistics", server.ValidateJwt(config), controller.GetStatistics)

	router.GET("reports", server.ValidateJwt(config), server.Role("admin"), controller.DownloadReport)
	router.PATCH(":id/complete", server.ValidateJwt(config), server.Role("admin"), controller.UpdateStatus(dao.OrderCompletedStatus))

}
