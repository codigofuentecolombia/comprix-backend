package dto

import (
	"comprix/app/domain/dao"
	"time"
)

type NewOrder struct {
	Products        []OrderProduct         `form:"products"         json:"products"         binding:"required,dive"`
	PaymentMethod   dao.OrderPaymentMethod `form:"payment_method"   json:"payment_method"   binding:"required"`
	ShippingAddress OrderShippingAddress   `form:"shipping_address" json:"shipping_address" binding:"required"`
	// Datos que se llenan
	Time         string
	ShippingCost float64
}

type OrderProduct struct {
	ID       uint `form:"id"       json:"id"       binding:"required"`
	Quantity uint `form:"quantity" json:"quantity" binding:"required"`
}

type OrderShippingAddress struct {
	Date        time.Time `form:"date"            json:"date"             binding:"required"`
	TimeID      string    `form:"time_id"         json:"time_id"             binding:"required"`
	Street      string    `form:"street"          json:"street"           binding:"required"`
	Reference   *string   `form:"reference"       json:"reference"        `
	PhoneNumber string    `form:"phone_number"    json:"phone_number"     binding:"required,max=20"`
	// City           string  `form:"city"            json:"city"             binding:"required"`
	// State          string  `form:"state"           json:"state"            binding:"required"`
	// Colony         string  `form:"colony"          json:"colony"           binding:"required"`
	// PostalCode     string  `form:"postal_code"     json:"postal_code"      binding:"required"`
	// ExternalNumber string  `form:"external_number" json:"external_number"  binding:"required"`
	// InternalNumber *string `form:"internal_number" json:"internal_number"`
}

type OrderDetail struct {
	User            dao.User                `json:"user"`
	Order           dao.Order               `json:"order"`
	Items           []dao.PageProduct       `json:"items"`
	ShippingAddress dao.UserShippingAddress `json:"shipping_address"`
}

type OrderStatistics struct {
	Total     float64 `json:"total"`
	Pending   int     `json:"pending"`
	Completed int     `json:"completed"`
}
