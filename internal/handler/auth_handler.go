package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/programmerjide/ecommerce/internal/dto"
	"github.com/programmerjide/ecommerce/internal/service"
	"github.com/programmerjide/ecommerce/internal/utils"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      zerolog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Registration failed")
		utils.BadRequestResponse(c, "Registration failed: ", err)
		return
	}

	utils.CreatedResponse(c, "User registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	h.logger.Info().Str("email", req.Email).Msg("Login attempt") // ✅ Add logging

	response, err := h.authService.Login(&req)
	if err != nil {
		h.logger.Error().Err(err).Str("email", req.Email).Msg("Login failed") // ✅ Log actual error
		utils.UnauthorizedResponse(c, err.Error())                            // ✅ Return actual error temporarily
		return
	}

	utils.CreatedResponse(c, "Login successful", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	response, err := h.authService.RefreshToken(&req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Token refresh failed")
		utils.UnauthorizedResponse(c, "Token refresh failed")
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		h.logger.Error().Err(err).Msg("Logout failed")
		utils.InternalServerErrorResponse(c, "Logout failed")
		return
	}

	utils.SuccessResponse(c, "Logout successful", nil)
}
