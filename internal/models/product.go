package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category in the e-commerce system
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Products []Product `json:"-" gorm:"foreignKey:CategoryID"` // ✅ Excluded
}

// Product represents a product in the e-commerce system
type Product struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CategoryID  uint           `json:"category_id" gorm:"not null"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Stock       int            `json:"stock" gorm:"default:0"`
	SKU         string         `json:"sku" gorm:"uniqueIndex;not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	IsOnSale    bool           `json:"is_on_sale" gorm:"default:true"`
	IsFeatured  bool           `json:"is_featured" gorm:"default:false"`
	Rating      float64        `json:"rating" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // ✅ Proper soft deletes

	Category   Category       `json:"category" gorm:"foreignKey:CategoryID"` // ✅ Included
	Images     []ProductImage `json:"images" gorm:"foreignKey:ProductID"`    // ✅ Included
	Tags       []string       `json:"tags" gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem    `json:"-" gorm:"foreignKey:ProductID"` // ✅ Excluded
	CartItems  []CartItem     `json:"-" gorm:"foreignKey:ProductID"` // ✅ Excluded
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID uint           `json:"product_id" gorm:"not null"`
	URL       string         `json:"url" gorm:"not null"`
	AltText   string         `json:"alt_text"`
	IsPrimary bool           `json:"is_primary" gorm:"default:false"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Product Product `json:"-" gorm:"foreignKey:ProductID"` // ✅ Excluded
}
