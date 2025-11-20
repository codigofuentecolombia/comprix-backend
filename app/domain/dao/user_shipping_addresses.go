package dao

import (
	"time"

	"gorm.io/gorm"
)

type UserShippingAddress struct {
	ID             uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	UserID         uint      `json:"user_id"         gorm:"not null;index"`
	Date           time.Time `json:"date"            gorm:"not null"`
	Time           string    `json:"time"            gorm:"not null"`
	Alias          string    `json:"alias"           gorm:"type:varchar(255);not null;default:'Principal'"`
	Street         string    `json:"street"          gorm:"type:varchar(255);not null"`
	Colony         string    `json:"colony"          gorm:"type:varchar(255);not null"`
	State          string    `json:"state"           gorm:"type:varchar(255);not null"`
	City           string    `json:"city"            gorm:"type:varchar(255);not null"`
	PostalCode     string    `json:"postal_code"     gorm:"type:varchar(20);not null"`
	PhoneNumber    string    `json:"phone_number"  gorm:"type:varchar(20);not null"`
	Reference      *string   `json:"reference"       gorm:"type:text"`
	ExternalNumber string    `json:"external_number" gorm:"type:varchar(50);not null"`
	InternalNumber *string   `json:"internal_number" gorm:"type:varchar(50)"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (UserShippingAddress) TableName() string {
	return "user_shipping_addresses"
}
