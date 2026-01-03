package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/programmerjide/ecommerce/internal/dto"
	"github.com/programmerjide/ecommerce/internal/service"
	"github.com/programmerjide/ecommerce/internal/utils"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	UserService *service.UserService
	logger      zerolog.Logger
}

func NewUserHandler(userService *service.UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		UserService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// Implementation for getting user profile
	userID := c.GetUint("user_id")
	profile, err := h.UserService.GetProfile(userID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get user profile")
		utils.NotFoundResponse(c, "User profile not found")
		return
	}

	utils.SuccessResponse(c, "User profile retrieved successfully", profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// Implementation for updating user profile
	userID := c.GetUint("user_id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request payload", err)
		return
	}

	updatedProfile, err := h.UserService.UpdateProfile(userID, &req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update user profile")
		utils.InternalServerErrorResponse(c, "Failed to update profile", err)
		return
	}

	utils.SuccessResponse(c, "User profile updated successfully", updatedProfile)
}
