package dto

type IPage interface {
	GetProductDetail(params PageParams) *RetrievedProduct
	GetAllProductsLinks(params PageParams) []string

	CreateOrUpdateProduct(product *RetrievedProduct)
	CreateOrUpdateProductsByLinks(links []string)
}

type PageParams struct {
	Url   string
	Tries int
}
