package order_controller

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	repository_page_product "comprix/app/repositories/page-product"
	"comprix/app/server"
	service_auth "comprix/app/services/auth"
	service_email "comprix/app/services/email"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type OrderController struct {
	email        service_email.Service
	config       dto.Config
	service      service_auth.Service
	repositories OrderControllerRepositories
}

type OrderControllerRepositories struct {
	user        repositories.UserRepository
	order       repositories.OrderRepository
	product     repositories.PageProductRepository
	calendar    repositories.CalendarRepository
	pageProduct *repository_page_product.Repository
}

func InitOrderController(config dto.Config) OrderController {
	return OrderController{
		email:   service_email.InitService(config),
		config:  config,
		service: service_auth.InitService(config),
		repositories: OrderControllerRepositories{
			user:        repositories.InitUserRepository(config.GormDB),
			order:       repositories.InitOrderRepository(config.GormDB),
			product:     repositories.InitPageProductRepository(config.GormDB),
			calendar:    repositories.InitCalendarRepository(config.GormDB),
			pageProduct: repository_page_product.InitRepository(config.GormDB),
		},
	}
}

func (ctr OrderController) New(ginContext *gin.Context) {
	// Obtener parametros validados
	form := server.GinFormBinding(ginContext, dto.NewOrder{}, "Invalid form")
	authUser := ginContext.MustGet("authUser").(*dao.User)
	// Obtener los productos
	var products []dao.PageProduct
	//
	for _, formProduct := range form.Products {
		currentPageProduct := ctr.repositories.pageProduct.Find(
			dto.ProductRepositoryFindParams{
				ID:       &formProduct.ID,
				Preloads: &[]dto.RepositoryGormParams{{Query: "Page"}, {Query: "Product"}},
			},
		)
		// Mandar error si ocurrio un error o el producto esta desactivado
		if currentPageProduct.Error != nil || currentPageProduct.Data.Product.IsDisabled {
			server.InternalErrorException("No se pudo obtener el producto", currentPageProduct.Error)
		}
		// Actualizar cantidad
		currentPageProduct.Data.Quantity = formProduct.Quantity
		// Agregar producto
		products = append(products, currentPageProduct.Data)
	}
	// Obtener fecha
	time, err := ctr.repositories.calendar.FindByID(form.ShippingAddress.TimeID)
	// Verificar si existe err
	if err != nil {
		server.InternalErrorException("No se pudo obtener el horario de entrega", nil)
	}
	// LLenar datos faltantes
	form.Time = fmt.Sprintf("%s - %s", time.StartTime, time.EndTime)
	form.ShippingCost = time.Price
	// Crear orden
	if order, err := ctr.repositories.order.Create(*authUser, products, form); err != nil {
		server.InternalErrorException("No se pudo obtener el producto", nil)
	} else {
		ctr.email.Order(*authUser, order, products)
	}
	// Regresar respuesta
	server.SuccessResponse(ginContext)
}

func (ctr OrderController) GetAll(ginContext *gin.Context) {
	claims := ginContext.MustGet("claims").(jwt.MapClaims)
	userID := claims["sub"].(map[string]interface{})["id"]
	// Obteenr orden
	orders, err := ctr.repositories.order.FindAllByUserID(userID)
	// Validar si se pudo obtener
	if err != nil {
		server.InternalErrorException("No se pudo obtener las ordenes", nil)
	}
	// Regresar
	server.Response(ginContext, http.StatusOK, orders)
}
