package dao

import (
	"time"
)

type ProductGroup struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement"`
	ReferenceProductID uint       `gorm:"not null"`
	CreatedAt          *time.Time `gorm:"column:created_at"`

	// Relaciones
	Members []ProductGroupMember `gorm:"foreignKey:GroupID"`
}

func (ProductGroup) TableName() string {
	return "product_groups"
}
