package dao

import (
	"time"

	"gorm.io/gorm"
)

type Brand struct {
	ID        int            `json:"id"        gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name"      gorm:"type:varchar(250);not null;unique"`
	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Brand) TableName() string {
	return "brands"
}
