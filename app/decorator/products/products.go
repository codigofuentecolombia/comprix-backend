package products_decorator

import "comprix/app/repositories"

type Decorator struct {
	repository repositories.PageProductRepository
}

func InitService() Decorator {
	return Decorator{}
}
