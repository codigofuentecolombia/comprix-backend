package dao

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	RoleID      uint8   `json:"-"             gorm:"type:int"`
	Email       string  `json:"email"         gorm:"type:varchar(250);not null;unique"`
	Picture     *string `json:"picture"       gorm:"type:varchar(250)"`
	Username    string  `json:"username"      gorm:"type:varchar(250);not null;unique"`
	Password    string  `json:"-"             gorm:"type:varchar(250);not null"`
	FirstName   string  `json:"first_name"    gorm:"type:varchar(250);not null"`
	LastName    string  `json:"last_name"     gorm:"type:varchar(250);not null"`
	PhoneNumber string  `json:"phone_number"  gorm:"type:varchar(20);not null"`
	IsVerified  bool    `json:"is_verified"   gorm:"type:int"`

	Role Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`

	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

// v=spf1 include:spf.hostmar.com -all
// 149.50.143.46
