package dao

import (
	"time"

	"gorm.io/gorm"
)

type Calendar struct {
	ID        uint    `json:"id"         gorm:"primaryKey;autoIncrement"`
	StartTime string  `json:"start_time" gorm:"type:time;not null"`
	EndTime   string  `json:"end_time"   gorm:"type:time;not null"`
	Price     float64 `json:"price"      gorm:"type:decimal(10,2)"`
	Monday    bool    `json:"monday"     gorm:"type:tinyint(1);not null;default:0"`
	Tuesday   bool    `json:"tuesday"    gorm:"type:tinyint(1);not null;default:0"`
	Wednesday bool    `json:"wednesday"  gorm:"type:tinyint(1);not null;default:0"`
	Thursday  bool    `json:"thursday"   gorm:"type:tinyint(1);not null;default:0"`
	Friday    bool    `json:"friday"     gorm:"type:tinyint(1);not null;default:0"`
	Saturday  bool    `json:"saturday"   gorm:"type:tinyint(1);not null;default:0"`
	Sunday    bool    `json:"sunday"     gorm:"type:tinyint(1);not null;default:0"`

	CreatedAt time.Time      `json:"created_at" gorm:"type:timestamp;not null;default:current_timestamp"`
	UpdatedAt *time.Time     `json:"-" gorm:"type:timestamp;default:null"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"type:timestamp;default:null" gorm:"index"`
}

func (Calendar) TableName() string {
	return "calendar"
}
