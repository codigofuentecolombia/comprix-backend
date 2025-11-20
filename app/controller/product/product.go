package product_controller

import (
	"comprix/app/domain/dto"
	"comprix/app/repositories"
	repository_category "comprix/app/repositories/category"
	repository_page_product "comprix/app/repositories/page-product"
	repository_product "comprix/app/repositories/product"
	pages_carrefour "comprix/app/scrapper/pages/carrefour"
	pages_hiperlibertad "comprix/app/scrapper/pages/hiperlibertad"
	pages_jumbo "comprix/app/scrapper/pages/jumbo"
	pages_masonline "comprix/app/scrapper/pages/masonline"
	pages_vea "comprix/app/scrapper/pages/vea"
)

type Controller struct {
	config       *dto.Config
	pages        *ControllerPages
	repositories *ControllerRepositories
}

type ControllerRepositories struct {
	err         *repositories.ErrorRepository
	product     *repository_product.Repository
	category    *repository_category.Repository
	pageProduct *repository_page_product.Repository
}

type ControllerPages struct {
	vea           dto.IScrapperService
	jumbo         dto.IScrapperService
	carrefour     dto.IScrapperService
	masonline     dto.IScrapperService
	hiperlibertad dto.IScrapperService
}

func InitController(config *dto.Config) Controller {
	vea, _ := pages_vea.Initialize(config)
	jumbo, _ := pages_jumbo.Initialize(config)
	carrefour, _ := pages_carrefour.Initialize(config)
	masonline, _ := pages_masonline.Initialize(config)
	hiperlibertad, _ := pages_hiperlibertad.Initialize(config)
	// Controllador
	return Controller{
		config: config,
		repositories: &ControllerRepositories{
			err:         repositories.InitErrorRepository(config.GormDB),
			product:     repository_product.InitRepository(config.GormDB),
			category:    repository_category.InitRepository(config.GormDB),
			pageProduct: repository_page_product.InitRepository(config.GormDB),
		},
		pages: &ControllerPages{
			vea:           vea,
			jumbo:         jumbo,
			carrefour:     carrefour,
			masonline:     masonline,
			hiperlibertad: hiperlibertad,
		},
	}
}
