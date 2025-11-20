package pages_alem

import "comprix/app/domain/dao"

func (service Service) GetExistingProducts() ([]dao.PageProduct, error) {
	return service.svc.GetExistingProducts()
}
