package pages_vea

import (
	"comprix/app/constants"
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
	service_scrapper "comprix/app/services/scrapper"
)

type Service struct {
	svc *service_scrapper.Service
}

func Initialize(conf *dto.Config) (dto.IScrapperService, error) {
	// Inicializar servicio de scrapper
	scrapper, err := service_scrapper.InitService(conf, constants.PageVea, 24)
	// Verificar si ocurrio error
	if err != nil {
		return nil, err
	}
	// Regresar
	return &Service{scrapper}, nil
}

func (service Service) CloseChromedpCtx() {
	service.svc.CloseChromedpCtx()
}

func (service Service) GetErrors() []dao.Error {
	return service.svc.GetErrors()
}
