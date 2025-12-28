package service

import (
	"errors"
	"time"

	"github.com/programmerjide/ecommerce/internal/config"
	"github.com/programmerjide/ecommerce/internal/dto"
	"github.com/programmerjide/ecommerce/internal/models"
	"github.com/programmerjide/ecommerce/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{
		db:     db,
		config: cfg,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	var existingUser models.User

	// ✅ FIXED: Check if user EXISTS (err == nil means found)
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already in use")
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := models.User{
		Email:     req.Email,
		Password:  hashPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Role:      models.UserRoleCustomer,
		IsActive:  true,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, errors.New("failed to create user")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	if err := s.db.Where("email = ? AND is_active = ?", req.Email, true).First(&user).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := utils.VerifyPassword(user.Password, req.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return s.generateAuthResponse(&user)
}

func (s *AuthService) RefreshToken(req *dto.RefreshTokenRequest) (*dto.AuthResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, s.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	var refreshToken models.RefreshToken
	if err := s.db.Where("token = ? AND expires_at > ?", req.RefreshToken, time.Now()).First(&refreshToken).Error; err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	var user models.User
	if err := s.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// Optional: Delete old refresh token before generating new one
	s.db.Delete(&refreshToken)

	return s.generateAuthResponse(&user)
}

func (s *AuthService) Logout(refreshToken string) error {
	result := s.db.Where("token = ?", refreshToken).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return errors.New("failed to logout")
	}
	if result.RowsAffected == 0 {
		return errors.New("refresh token not found")
	}
	return nil
}

func (s *AuthService) generateAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, refreshToken, err := utils.GenerateJWTToken(
		&s.config.JWT,
		user.ID,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, errors.New("failed to generate tokens")
	}

	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.config.JWT.RefreshTokenExpires),
	}

	if err := s.db.Create(refreshTokenModel).Error; err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Phone:     user.Phone,
			Role:      string(user.Role),
			IsActive:  user.IsActive,
		},
	}, nil
}
