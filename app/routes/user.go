package routes

import (
	user_controller "comprix/app/controller/user"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func UserRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/users")

	controller := user_controller.InitController(config)

	router.GET("paginated", server.ValidateJwt(config), server.Role("admin"), controller.GetPaginated)
}
