package service_scrapper

import (
	"comprix/app/domain/dao"
	"comprix/app/domain/dto"
)

func (svc *Service) GetErrors() []dao.Error {
	params := dto.ErrorRepositoryFindParams{PageID: &svc.Page.ID}
	response := svc.Repositories.Error.FindAll(params)
	//
	return response.Data
}
