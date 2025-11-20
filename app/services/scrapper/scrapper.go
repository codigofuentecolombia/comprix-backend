package service_scrapper

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	"comprix/app/fails"
	"comprix/app/logger"
	"comprix/app/repositories"
	service_product "comprix/app/services/product"
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type Service struct {
	Conf         *dto.Config
	Logs         *ServiceLogs
	Page         *dao.Page
	ChromedpCtx  *ServiceChromedpCtx
	Repositories *ServiceRepositories
	//
	InitializedAt   time.Time
	ProductsPerPage float64
	//
	productSvc *service_product.Service
}

type ServiceChromedpCtx struct {
	Ctx       context.Context
	CancelFns ServiceChromedpCtxCancelFn
}

type ServiceChromedpCtxCancelFn struct {
	Chromedp  context.CancelFunc
	Allocator context.CancelFunc
}

type ServiceLogs struct {
	Debug   *logrus.Logger
	Warning *logrus.Logger
}

type ServiceRepositories struct {
	Error   *repositories.ErrorRepository
	Product *repositories.ProductRepository
}

func InitService(conf *dto.Config, pageName string, productsPerPage float64) (*Service, error) {
	svc := &Service{
		Conf:            conf,
		productSvc:      service_product.InitService(conf),
		InitializedAt:   time.Now(),
		ProductsPerPage: productsPerPage,
		//
		Logs: &ServiceLogs{
			Debug:   logger.CreateDebug(conf.Settings.Paths.Logs, pageName, conf.Settings.Server.Debug),
			Warning: logger.CreateWarning(conf.Settings.Paths.Logs, pageName, conf.Settings.Server.Debug),
		},
		Repositories: &ServiceRepositories{
			Error:   repositories.InitErrorRepository(conf.GormDB),
			Product: repositories.InitProductRepository(conf.GormDB),
		},
	}
	// Obtener pagina
	page, err := repositories.InitPageRepository(conf.GormDB).FindByName(pageName)
	// Verificar si existe error
	if err != nil {
		return nil, fails.Create("No se pudo obtener la pagina desde la base de datos", err)
	}
	// Actualizar elemento de pagina
	svc.Page = page
	// // Obtener contexto
	// ctx, cancelFns, err := svc.InitChromedp()
	// // Regresar error
	// if err != nil {
	// 	return nil, fails.Create("No se pudo inicializar chromedp", err)
	// }
	// // Actualizar
	// svc.ChromedpCtx = &ServiceChromedpCtx{
	// 	Ctx:       ctx,
	// 	CancelFns: cancelFns,
	// }
	// Regresar sin error
	return svc, nil
}
