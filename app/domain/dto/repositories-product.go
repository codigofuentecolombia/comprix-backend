package dto

import "comprix/app/domain/dao"

type ProductRepositoryFindParams struct {
	ID                   *uint `json:"id" uri:"id"`
	Url                  *string
	Name                 *string
	ExcludeID            *uint
	ProductID            *uint
	BestPrice            *bool
	Softdeleted          *bool
	OmitDisabled         *bool
	OnlyDisabled         *bool
	BestPagePrice        *bool
	OlderThanOneDay      *bool
	OnlyRecommended      *bool
	OnlyWithDiscount     *bool
	WithDistinctProducts *bool
	// Query params
	Type       *PageProductType `json:"-" form:"type,omitempty"`
	Order      *string          `json:"-" form:"order_by,omitempty"`
	Search     *string          `json:"-" form:"search,omitempty"`
	PageID     *uint            `json:"-" form:"page_id,omitempty"`
	BranchIDS  *[]uint          `json:"-" form:"branch_ids,omitempty"`
	CategoryID *string          `json:"-" form:"category,omitempty"`
	//
	Selects  *RepositoryGormSelections
	OrderBy  any
	Preloads *[]RepositoryGormParams

	Pagination *Pagination[dao.PageProduct] `form:"pagination, dive"`
}

type PageProductType string

const AllPageProductType PageProductType = ""
const OffersPageProductType PageProductType = "offers"
const DisablePageProductType PageProductType = "disables"
const DiscountsPageProductType PageProductType = "discounts"
const RecommendedPageProductType PageProductType = "recommended"

type ProductStatusColumn string

const ProductDisableStatusColumn ProductStatusColumn = "is_disabled"
const ProductDiscountStatusColumn ProductStatusColumn = "is_in_discount"
const ProductRecommendedStatusColumn ProductStatusColumn = "is_recommended"
