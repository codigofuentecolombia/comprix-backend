package service_product

import (
	"comprix/app/domain/dto"
	repository_category "comprix/app/repositories/category"
	repository_page_product "comprix/app/repositories/page-product"
	repository_product "comprix/app/repositories/product"
)

type Service struct {
	config       *dto.Config
	repositories *ServiceRepositories
}

type ServiceRepositories struct {
	product     *repository_product.Repository
	category    *repository_category.Repository
	pageProduct *repository_page_product.Repository
}

func InitService(cnf *dto.Config) *Service {
	return &Service{
		config: cnf,
		repositories: &ServiceRepositories{
			product:     repository_product.InitRepository(cnf.GormDB),
			category:    repository_category.InitRepository(cnf.GormDB),
			pageProduct: repository_page_product.InitRepository(cnf.GormDB),
		},
	}
}
