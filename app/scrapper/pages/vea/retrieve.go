package pages_vea

import "comprix/app/domain/dao"

func (service Service) GetExistingProducts() ([]dao.PageProduct, error) {
	return service.svc.GetExistingProducts()
}
