package pages_vea

import (
	"comprix/app/domain/dto"
)

func (service Service) CreateOrUpdateProduct(pageProduct *dto.RetrievedProduct) {
	service.svc.HandleProduct(pageProduct)
}

func (service Service) DisableNotFound() {
	service.svc.DisableNotFound()
}

func (service Service) GetNotFoundProductLinks() []string {
	return service.svc.GetNotFoundProductLinks()
}
