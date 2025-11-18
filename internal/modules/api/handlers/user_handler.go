package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"your-project/internal/dto"
	"your-project/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetCurrentUser godoc
// @Summary Get current user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]interface{}
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.GetString("userID")

	user, err := h.userService.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found", "code": 404})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser godoc
// @Summary Update current user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "Update data"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]interface{}
// @Router /users/me [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID := c.GetString("userID")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	user, err := h.userService.UpdateUser(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteCurrentUser godoc
// @Summary Delete current user
// @Tags users
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} map[string]interface{}
// @Router /users/me [delete]
func (h *UserHandler) DeleteCurrentUser(c *gin.Context) {
	userID := c.GetString("userID")

	err := h.userService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.Status(http.StatusNoContent)
}

// InviteUser godoc
// @Summary Invite team member
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.InviteUserRequest true "Invitation data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/invite [post]
func (h *UserHandler) InviteUser(c *gin.Context) {
	var req dto.InviteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	err := h.userService.InviteUser(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "code": 400})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent"})
}
