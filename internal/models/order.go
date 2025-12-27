// Description: This file defines the Order, OrderItem, Cart, and CartItem models for an e-commerce system using GORM.
package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a customer's order in the e-commerce system
type Order struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Status      OrderStatus    `json:"status" gorm:"default:'pending'"`
	TotalAmount float64        `json:"total_amount" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	User       User        `json:"user" gorm:"foreignKey:UserID"`         // ✅ Included
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"` // ✅ Included
}

// OrderStatus defines the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // default status
	OrderStatusPaid      OrderStatus = "paid"      // payment received
	OrderStatusShipped   OrderStatus = "shipped"   // order shipped
	OrderStatusDelivered OrderStatus = "delivered" // order delivered
	OrderStatusCancelled OrderStatus = "cancelled" // order cancelled
	OrderStatusRefunded  OrderStatus = "refunded"  // order refunded
	OrderStatusConfirmed OrderStatus = "confirmed" // order confirmed
	OrderStatusFailed    OrderStatus = "failed"    // payment failed
)

// OrderItem represents an item in a customer's order
type OrderItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	OrderID   uint           `json:"order_id" gorm:"not null"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	Price     float64        `json:"unit_price" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Order   Order   `json:"-" gorm:"foreignKey:OrderID"`         // ✅ Excluded
	Product Product `json:"product" gorm:"foreignKey:ProductID"` // ✅ Included
}

// Cart represents a shopping cart for a user
type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null;uniqueIndex"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationship
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"` // ✅ Included
}

// CartItem represents an item in a shopping cart
type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CartID    uint           `json:"cart_id" gorm:"not null"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	Quantity  int            `json:"quantity" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Cart    Cart    `json:"-" gorm:"foreignKey:CartID"`          // ✅ Excluded
	Product Product `json:"product" gorm:"foreignKey:ProductID"` // ✅ Included
}
