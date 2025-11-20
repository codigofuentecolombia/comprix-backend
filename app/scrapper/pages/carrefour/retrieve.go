package pages_carrefour

import "comprix/app/domain/dao"

func (service Service) GetExistingProducts() ([]dao.PageProduct, error) {
	return service.svc.GetExistingProducts()
}
