package routes

import (
	"comprix/app/domain/dto"

	"github.com/gin-gonic/gin"
)

func RouterHandler(config dto.Config, ginRouter *gin.Engine) {
	// Crear router
	router := ginRouter.Group("/api/v1")
	// Rutas
	HealthRoutes(config, router)
	AuthRoutes(config, router)
	UserRoutes(config, router)
	PageRoutes(config, router)
	ErrorRoutes(config, router)
	OrderRoutes(config, router)
	BrandsRoutes(config, router)
	ProductRoutes(config, router)
	CategoryRoutes(config, router)
	CalendarRoutes(config, router)
	AdminProductRoutes(config, router)
	ScrapperRoutes(config, router)
}
