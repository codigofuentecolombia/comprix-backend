package dao

import (
	"time"

	"gorm.io/gorm"
)

type OrderProduct struct {
	ID            uint64          `gorm:"primaryKey;autoIncrement"`
	Quantity      uint            `gorm:"not null;index"`
	OrderID       uint            `gorm:"not null;index"`
	PageProductID uint            `gorm:"not null;index"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime"`
	DeletedAt     *gorm.DeletedAt `gorm:"index"`
	// Relación con PageProduct
	PageProduct PageProduct `gorm:"foreignKey:PageProductID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}
