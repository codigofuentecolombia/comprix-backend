package dto

type BrandRepositoryFindParams struct {
	OrderBy     *BrandRepositoryFindOrderBy
	HasProducts *bool
}

type BrandRepositoryFindOrderBy struct {
	NameDesc          *bool
	MaxNameLengthDesc *bool
}
