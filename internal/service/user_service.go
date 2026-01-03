package service

import (
	"github.com/programmerjide/ecommerce/internal/dto"
	"github.com/programmerjide/ecommerce/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) GetProfile(userID uint) (*dto.UserResponse, error) {
	// Implement logic to retrieve user profile from the database
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
	}, nil
}

func (s *UserService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*dto.UserResponse, error) {
	// Implement logic to update user profile in the database
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return s.GetProfile(userID) // Return the updated profile
}
