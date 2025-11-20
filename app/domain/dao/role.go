package dao

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID uint8 `json:"id" gorm:"primaryKey;autoIncrement"`

	Name        string `json:"name"         gorm:"type:varchar(250);not null;unique"`
	DisplayName string `json:"display_name" gorm:"type:varchar(250);not null"`
	Description string `json:"description"  gorm:"type:text;not null;default:''"`

	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Role) TableName() string {
	return "roles"
}
