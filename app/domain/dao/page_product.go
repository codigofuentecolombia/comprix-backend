package dao

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type PageProduct struct {
	ID            uint     `json:"id"              gorm:"primaryKey;autoIncrement"`
	Url           string   `json:"url,omitempty"             gorm:"nol null"`
	Images        []string `json:"images,omitempty"          gorm:"serializer:json"`
	PageID        uint     `json:"page_id,omitempty"         gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductID     uint     `json:"product_id,omitempty"      gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MainProductID uint     `json:"main_product_id,omitempty" gorm:"not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	Price                 float64 `json:"price"                             gorm:"type:decimal(10,2)"`
	DiscountPrice         float64 `json:"discount_price"                    gorm:"type:decimal(10,2)"`
	OriginalPrice         string  `json:"original_price,omitempty"          gorm:"type:varchar(100)"`
	OriginalDiscountPrice string  `json:"original_discount_price,omitempty" gorm:"type:varchar(100)"`

	MinQuantityToApplyDiscount uint `json:"min_quantity_to_apply_discount" gorm:"type:int(11)"`

	CreatedAt *time.Time            `json:"created_at"`
	UpdatedAt *time.Time            `json:"-"`
	DeletedAt soft_delete.DeletedAt `json:"-" gorm:"index"`

	Quantity uint `json:"quantity,omitempty" gorm:"->;<-:false"`

	Page    *Page    `json:"page,omitempty" gorm:"foreignKey:PageID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:MainProductID"`

	OtherStores *[]PageProduct `json:"other_stores,omitempty" gorm:"-"`
}

func (PageProduct) TableName() string {
	return "page_products"
}
