package controllers

import (
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userUseCase usecases.UserUseCase
}

func NewUserController(userUseCase usecases.UserUseCase) *UserController {
	return &UserController{userUseCase: userUseCase}
}

func (c *UserController) GetAllUsers(ctx *gin.Context) {
	roleInterface, exists := ctx.Get("role")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Role tidak ditemukan"})
		return
	}

	role, ok := roleInterface.(string)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role type"})
		return
	}

	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user_id type"})
		return
	}
	users, err := c.userUseCase.GetAllUsers(role, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve users",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data":    users,
	})
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID format",
		})
		return
	}
	user, err := c.userUseCase.GetUserByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "User not found",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User retrieved successfully",
		"data":    user,
	})
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var request entities.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	roleInterface, _ := ctx.Get("role")
	role, _ := roleInterface.(string)

	userIDInterface, _ := ctx.Get("user_id")
	var userID uint
	if userIDFloat, ok := userIDInterface.(float64); ok {
		userID = uint(userIDFloat)
	} else if userIDUint, ok := userIDInterface.(uint); ok {
		userID = userIDUint
	}

	var (
		user entities.UserResponse
		err  error
	)

	if role == "admin" {
		user, err = c.userUseCase.CreateFromAdmin(request)
	} else {
		user, err = c.userUseCase.CreateFromManajement(userID, request)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data":    user,
	})
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID format",
		})
		return
	}
	var request entities.CreateUserRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	user, err := c.userUseCase.UpdateUser(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
		"data":    user,
	})
}

func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid user ID format",
		})
		return
	}

	err = c.userUseCase.DeleteUser(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete user",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}

func (c *UserController) GetMe(ctx *gin.Context) {
	userIDInterface, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User ID not found",
		})
		return
	}

	var userID uint
	if userIDFloat, ok := userIDInterface.(float64); ok {
		userID = uint(userIDFloat)
	} else if userIDUint, ok := userIDInterface.(uint); ok {
		userID = userIDUint
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid user ID type",
		})
		return
	}

	user, err := c.userUseCase.GetUserByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve user data",
			"error":   err.Error(),
		})
		return
	}

	response := entities.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User data retrieved successfully",
		"data":    response,
	})
}
