package pages_alem

import (
	"comprix/app/domain/dto"
)

func (service Service) CreateOrUpdateProduct(pageProduct *dto.RetrievedProduct) {
	service.svc.HandleProduct(pageProduct)
}

func (service Service) DisableNotFound() {
	service.svc.DisableNotFound()
}
