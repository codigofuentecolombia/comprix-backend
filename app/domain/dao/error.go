package dao

import (
	"time"
)

type Error struct {
	ID         uint      `json:"id"          gorm:"primaryKey;autoIncrement"`
	Type       string    `json:"-"           gorm:"type:varchar(50);not null;default:'scraping'"`
	Url        string    `json:"url"         gorm:"type:varchar(255);not null"`
	Error      string    `json:"ctx"         gorm:"type:text;not null"`
	PageID     uint      `json:"-"           gorm:"not null"`
	Message    string    `json:"message"     gorm:"type:varchar(255);not null"`
	CreatedAt  time.Time `json:"created_at"  gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"-"           gorm:"autoUpdateTime"`
	Categories []string  `json:"categories"  gorm:"serializer:json"`
	// Relaciones
	Page *Page `json:"page,omitempty" gorm:"foreignKey:PageID"`
}

func (Error) TableName() string {
	return "errors"
}
