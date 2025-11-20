package dto

type RetrievedProducts []RetrievedProduct

type RetrievedProduct struct {
	Url                   string   `json:"url"`
	Sku                   string   `json:"sku"`
	Name                  string   `json:"name"`
	Price                 float64  `json:"price"`
	Brand                 string   `json:"brand"`
	PageID                uint     `json:"page_id"`
	Images                []string `json:"images"`
	HasStock              *bool    `json:"has_stock"`
	Categories            []string `json:"categories"`
	Description           string   `json:"description"`
	DiscountPrice         float64  `json:"discount_price"`
	OriginalPrice         string   `json:"original_price"`
	OriginalDiscountPrice string   `json:"original_discount_price"`
	// Datos que se llenan manualmente
	ProductID *uint
	//
	MinQuantityToApplyDiscount uint `json:"min_quantity_to_apply_discount"`
}

type UpdatePageProduct struct {
	ID uint `json:"id" uri:"id" binding:"required"`
	// Cuerpo
	Name          string  `json:"name"           form:"name"           binding:"required"`
	Price         float64 `json:"price"          form:"price"          binding:"required"`
	Description   string  `json:"description"    form:"description"    `
	DiscountPrice float64 `json:"discount_price" form:"discount_price" `
	//
	MinQuantityToApplyDiscount uint `json:"min_quantity_to_apply_discount" form:"min_quantity_to_apply_discount" binding:"required"`
}

type RetrievedProductPrice struct {
	Price         string `json:"price"`
	MinQuantity   uint   `json:"min_quantity"`
	DiscountPrice string `json:"discount_price"`
}

type GroupedProducts map[uint][]uint

type ScanPageProduct struct {
	Url         string `json:"url"         form:"url"         binding:"required"`
	Product     string `json:"product"     form:"product"     binding:"required"`
	Category    string `json:"category"    form:"category"    binding:"required"`
	Subcategory string `json:"subcategory" form:"subcategory" binding:"required"`
}

type GroupPageProduct struct {
	ID    uint `json:"id"     uri:"id"     binding:"required"`
	NewID uint `json:"new_id" uri:"new_id" binding:"required"`
}
