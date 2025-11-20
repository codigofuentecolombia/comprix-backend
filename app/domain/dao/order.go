package dao

import (
	"time"

	"gorm.io/gorm"
)

type OrderPaymentMethod string

const OrderCashPaymentMethod OrderPaymentMethod = "cash"
const OrderMercadoPagoPaymentMethod OrderPaymentMethod = "mercado_pago"
const OrderBankTransferPaymentMethod OrderPaymentMethod = "bank_transfer"

type OrderStatus string

const OrderPendingStatus OrderStatus = "pending"
const OrderCompletedStatus OrderStatus = "completed"
const OrderCancelledStatus OrderStatus = "cancelled"
const OrderProcessingStatus OrderStatus = "processing"

type Order struct {
	ID     uint        `json:"id"                   gorm:"primaryKey;autoIncrement"`
	UserID uint        `json:"-"                    gorm:"not null;index"`
	Status OrderStatus `json:"status"               gorm:"type:enum('pending','processing','completed','cancelled');default:'pending'"`

	Total         float64 `json:"total"          gorm:"not null"`
	Subtotal      float64 `json:"subtotal"       gorm:"not null"`
	TotalDiscount float64 `json:"total_discount" gorm:"not null"`
	ShippingCost  float64 `json:"shipping_cost"  gorm:"not null"`

	PaymentMethod         OrderPaymentMethod `json:"payment_method"       gorm:"not null"`
	UserShippingAddressID uint               `json:"-"                    gorm:"not null;index"`
	CreatedAt             time.Time          `json:"created_at"           gorm:"autoCreateTime"`
	UpdatedAt             time.Time          `json:"-"                    gorm:"autoUpdateTime"`
	DeletedAt             *gorm.DeletedAt    `json:"-"                    gorm:"index"`

	// Relaciones
	User                User                `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	UserShippingAddress UserShippingAddress `json:"-" gorm:"foreignKey:UserShippingAddressID;constraint:OnDelete:CASCADE"`
	// Mostrar en json
	Items []OrderProduct `json:"items" gorm:"foreignKey:OrderID"`
}
