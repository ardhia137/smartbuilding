package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smartbuilding/entities"
	"smartbuilding/usecases"
)

type AuthController struct {
	authUseCase usecases.AuthUseCase
}

func NewAuthController(authUseCase usecases.AuthUseCase) *AuthController {
	return &AuthController{authUseCase: authUseCase}
}
func (c *AuthController) Login(ctx *gin.Context) {
	var request entities.LoginRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	response, err := c.authUseCase.Login(request.Email, request.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Failed to login",
			"error":   err.Error(),
		})
		return
	}

	ctx.Set("role", response.Role)
	ctx.Set("user_id", response.UserId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User logged in successfully",
		"data":    response,
	})
}
func (c *AuthController) ValidateToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
		return
	}

	user, err := c.authUseCase.ValidateToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid or expired token",
			"error":   err.Error(),
		})
		return
	}

	ctx.Set("role", user.Role)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token is valid",
		"data":    user,
	})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
		return
	}

	response, err := c.authUseCase.RefreshToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Failed to refresh token",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token refreshed successfully",
		"data":    response,
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
		return
	}

	err := c.authUseCase.Logout(token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to logout",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User logged out successfully",
	})
}
