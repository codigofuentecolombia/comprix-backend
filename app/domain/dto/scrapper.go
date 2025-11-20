package dto

import "comprix/app/domain/dao"

type ScrapperParams struct {
	Url        string
	Page       int
	Tries      int
	Categories []string
}

type ScrapperCategoryLink struct {
	Link       string
	Categories []string
}

type CreateOrUpdateScrappedProductsParams struct {
	Links      []string
	Categories []string
}

type ScrappedProducts []ScrappedProduct

type ScrappedProduct struct {
	Details         map[string]string `json:"details,omitempty"` // Detalles como un mapa
	Breadcrumbs     string            `json:"breadcrumbs"`
	AdditionalInfo  map[string]string `json:"additional_information,omitempty"` // Información adicional como un mapa
	Sku             string            `json:"sku"`
	Name            string            `json:"name"`
	Brand           string            `json:"brand"`
	Images          []string          `json:"images"`
	Category        string            `json:"category"`
	Description     string            `json:"description"`
	Price           string            `json:"price"`
	IncludedTaxes   string            `json:"included_taxes"`
	PriceCurrency   string            `json:"price_currency"`
	PriceValidUntil string            `json:"price_valid_until"`
}

type IScrapperService interface {
	//
	CloseChromedpCtx()
	//
	DisableNotFound()
	//
	GetExistingProducts() ([]dao.PageProduct, error)
	// Obtener links de categorias
	GetCategoryLinks() []ScrapperCategoryLink
	// Obtener detalle de producto
	GetProductDetail(params ScrapperParams) *RetrievedProduct
	// Obtener links de productos por pagina
	GetProductLinksByPage(params ScrapperParams) []string
	// Obtener total de paginas
	GetTotalPages(params ScrapperParams) int
	// Crear productos por categoria
	CreateOrUpdateProduct(pageProduct *RetrievedProduct)
	// Obtener urls con error
	GetErrors() []dao.Error
	// Obtener productos no encontrados
	GetNotFoundProductLinks() []string
}

type InitScrapPageFn func(conf *Config) (IScrapperService, error)
