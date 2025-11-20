package dao

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint32     `json:"id"            gorm:"primaryKey;autoIncrement"`
	Name          string     `json:"name"          gorm:"type:varchar(250);not null;unique"`
	ParentID      *uint32    `json:"parent_id"     gorm:"type:int"`
	Subcategories []Category `json:"subcategories" gorm:"foreignKey:ParentID;references:ID"`

	CreatedAt *time.Time     `json:"-"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Category) TableName() string {
	return "categories"
}
