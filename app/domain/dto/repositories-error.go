package dto

import "comprix/app/domain/dao"

type ErrorRepositoryFindParams struct {
	Url             *string
	PageID          *uint `json:"pageID" form:"pageID"`
	ShouldOrderDesc *bool
	//
	Preloads *[]RepositoryGormParams
	//
	Pagination Pagination[dao.Error]
}
