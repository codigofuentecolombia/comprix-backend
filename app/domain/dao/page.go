package dao

import (
	"time"

	"gorm.io/gorm"
)

type Page struct {
	ID   uint   `json:"id"   gorm:"primaryKey;autoIncrement;type:bigint"`
	Url  string `json:"-"    gorm:"type:varchar(255);not null;unique"`
	Logo string `json:"logo" gorm:"type:varchar(255);not null"`
	Name string `json:"name" gorm:"type:varchar(255);not null"`

	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
