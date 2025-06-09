package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"smartbuilding/entities"
	"smartbuilding/usecases"
)

type SettingController struct {
	haosUseCase usecases.SettingUseCase
}

func NewSettingController(haosUseCase usecases.SettingUseCase) *SettingController {
	return &SettingController{haosUseCase: haosUseCase}
}

func (c *SettingController) CreateSetting(ctx *gin.Context) {
	var request entities.CreateSettingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.haosUseCase.CreateSetting(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Setting created successfully",
		"data":    response,
	})
}

func (c *SettingController) GetAllSetting(ctx *gin.Context) {
	// Ambil role dari context
	roleInterface, _ := ctx.Get("role")
	role, _ := roleInterface.(string)

	// Ambil user_id dari context
	userIDInterface, _ := ctx.Get("user_id")
	userID, _ := userIDInterface.(uint)
	fmt.Println(userID)
	response, err := c.haosUseCase.GetAllSetting(role, userID)

	if err != nil {
		// Jika error "unauthorized", kirim status 401
		if err.Error() == "no data" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No Data Available"})
			return
		}
		// Error lain -> status 500
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Settings retrieved successfully",
		"data":    response,
	})
}

func (c *SettingController) GetSettingByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := c.haosUseCase.GetSettingByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Setting retrieved successfully",
		"data":    response,
	})
}

func (c *SettingController) UpdateSetting(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var request entities.CreateSettingRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.haosUseCase.UpdateSetting(id, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Setting updated successfully",
		"data":    response,
	})
}

func (c *SettingController) DeleteSetting(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.haosUseCase.DeleteSetting(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Setting deleted successfully",
	})
}
