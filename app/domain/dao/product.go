package dao

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint   `json:"id"             gorm:"primaryKey;autoIncrement;type:bigint"`
	Sku           string `json:"sku"            gorm:"type:varchar(100);not null;unique"`
	Name          string `json:"name"           gorm:"type:varchar(255);not null"`
	BrandID       int    `json:"brand_id"       gorm:"type:int"`
	CategoryID    uint32 `json:"category_id"    gorm:"type:int"`
	IsDisabled    bool   `json:"is_disabled"    gorm:"column:is_disabled"`
	Description   string `json:"description"    gorm:"type:text"`
	IsInDiscount  bool   `json:"is_in_discount" gorm:"column:is_in_discount"`
	IsRecommended bool   `json:"is_recommended" gorm:"column:is_recommended"`

	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Brand    Brand    `json:"brand,omitempty"    gorm:"foreignKey:BrandID"`
	Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Product) TableName() string {
	return "products"
}
