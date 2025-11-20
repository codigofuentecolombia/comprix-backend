package dto

import (
	"comprix/app/domain/dao"

	"gorm.io/gorm"
)

type RepositoryFindParams struct {
	Db         *gorm.DB
	Table      string
	Model      any
	Joins      []RepositoryFindCondition
	Preload    []RepositoryFindCondition
	GroupBy    *string
	Conditions RepositoryFindCondition
}

type RepositoryFindCondition struct {
	Query  string
	Params []interface{}
}

type ProductsType string

const AllProductsType ProductsType = ""
const OffersProductsType ProductsType = "offers"
const DiscountsProductsType ProductsType = "discounts"
const RecommendedProductsType ProductsType = "recommended"

type GetProductsParams struct {
	Type       ProductsType
	OrderBy    *string
	BranchIds  *[]string
	CategoryID string
	Pagination Pagination[dao.PageProduct]
}

type GetOrdersParams struct {
	UserID     any
	Pagination Pagination[dao.Order]
}

type GetUsersParams struct {
	Pagination Pagination[dao.User]
}

type RepositoryGenericResponse[T any] struct {
	Data  T
	Error error
}

type RepositoryGormParams struct {
	Args  any
	Query string
}

type RepositoryGormSelections struct {
	Args  []interface{}
	Query interface{}
}
