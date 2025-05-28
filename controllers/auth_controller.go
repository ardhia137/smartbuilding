package controllers

import (
	"fmt"
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
// @Summary Login pengguna
// @Description Melakukan autentikasi pengguna dengan email dan password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body entities.LoginRequest true "Kredensial login"
// @Success 200 {object} map[string]interface{} "Login berhasil"
// @Failure 400 {object} map[string]interface{} "Invalid input data"
// @Failure 401 {object} map[string]interface{} "Gagal login"
// @Router /auth/login [post]
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
// @Summary Validasi token
// @Description Memvalidasi token JWT
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Token valid"
// @Failure 400 {object} map[string]interface{} "Token diperlukan"
// @Failure 401 {object} map[string]interface{} "Token tidak valid atau kedaluwarsa"
// @Router /auth/validate [get]
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

// @Summary Memperbaharui token
// @Description Memperbaharui token JWT yang akan kedaluwarsa
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Token berhasil diperbaharui"
// @Failure 400 {object} map[string]interface{} "Token diperlukan"
// @Failure 401 {object} map[string]interface{} "Gagal memperbaharui token"
// @Router /auth/refresh [post]
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

// @Summary Logout pengguna
// @Description Menghapus token JWT dari sistem
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Berhasil logout"
// @Failure 400 {object} map[string]interface{} "Token diperlukan"
// @Failure 500 {object} map[string]interface{} "Gagal logout"
// @Router /auth/logout [post]
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
// @Summary Mengubah password
// @Description Mengubah password pengguna yang sedang login
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body entities.ChangePasswordRequest true "Data password lama dan baru"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Password berhasil diubah"
// @Failure 400 {object} map[string]interface{} "Token diperlukan atau request tidak valid"
// @Failure 401 {object} map[string]interface{} "Gagal mengubah password"
// @Router /auth/change-password [post]
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Token is required",
		})
		return
	}

	var req entities.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	err := c.authUseCase.ChangePassword(token, req.OldPassword, req.NewPassword)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Failed to change password",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Password changed successfully",
	})
}
