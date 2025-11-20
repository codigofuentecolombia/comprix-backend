package routes

import (
	"comprix/app/controller"
	product_controller "comprix/app/controller/product"
	"comprix/app/domain/dto"
	"comprix/app/server"

	"github.com/gin-gonic/gin"
)

func ProductRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/products")
	controllerV3 := product_controller.InitController(&config)
	controllerV1 := controller.InitProductController(config)
	//V3
	router.GET("discounts", controllerV3.GetWithDiscount)
	router.GET("paginated", controllerV3.GetPaginated)
	router.GET("recommended", controllerV3.GetRecommended)
	// Listrar rutas
	router.GET("", controllerV1.GetAll)
	router.GET("/:id", controllerV3.GetCheapest(true))
	router.GET("/outstanding", controllerV1.GetOutstanding)

	router.POST("/:id/recommended", server.ValidateJwt(config), server.Role("admin"), controllerV1.SetAsRecommended)
	router.DELETE("/:id/recommended", server.ValidateJwt(config), server.Role("admin"), controllerV1.UnsetAsRecommended)

	router.POST("/:id/in-discount", server.ValidateJwt(config), server.Role("admin"), controllerV1.SetAsInDiscount)
	router.DELETE("/:id/in-discount", server.ValidateJwt(config), server.Role("admin"), controllerV1.UnsetAsInDiscount)
}

func AdminProductRoutes(config dto.Config, ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("admin/products", server.ValidateJwt(config), server.Role("admin"))
	controller := controller.InitPageProductController(&config)
	controllerV2 := product_controller.InitController(&config)
	//
	router.GET("disables", controllerV2.GetDisabled)
	// Rutas
	router.GET("", controllerV2.GetAll)
	router.PUT("/:id", controller.Update)
	router.GET("/:id", controllerV2.GetCheapest(false))
	router.PUT("/:id/group/:new_id", controllerV2.GroupExisting)
	// Actualizar status
	router.POST("/:id/disable", controllerV2.UpdateStatus("is_disabled", true))
	router.DELETE("/:id/disable", controllerV2.UpdateStatus("is_disabled", false))

	router.POST("/scan", controllerV2.Scan)
}
