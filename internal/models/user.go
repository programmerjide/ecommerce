package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the e-commerce system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Phone     string         `json:"phone"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	Role      UserRole       `json:"role" gorm:"type:varchar(20);default:customer"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"` // ✅ Proper soft deletes

	RefreshTokens []RefreshToken `json:"-" gorm:"foreignKey:UserID"`
	Orders        []Order        `json:"-" gorm:"foreignKey:UserID"`
	Cart          Cart           `json:"-" gorm:"foreignKey:UserID"`
}

// UserRole defines the role of a user in the system
type UserRole string

const (
	UserRoleCustomer UserRole = "customer" // default role
	UserRoleAdmin    UserRole = "admin"    // admin role
)

// RefreshToken represents a refresh token for user authentication
type RefreshToken struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	Token     string         `json:"token" gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"-"`
}
