package controllers

import (
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MonitoringDataController struct {
	useCase    usecases.MonitoringDataUseCase
	hakAksesUc usecases.HakAksesUseCase
}

func NewMonitoringDataController(useCase usecases.MonitoringDataUseCase, hakAksesUc usecases.HakAksesUseCase) *MonitoringDataController {
	return &MonitoringDataController{useCase, hakAksesUc}
}

func (c *MonitoringDataController) SaveMonitoringData(ctx *gin.Context) {
	var requestData entities.CreateMonitoringDataRequest
	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	response, err := c.useCase.SaveMonitoringData(requestData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *MonitoringDataController) GetAirMonitoringData(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	roleInterface, _ := ctx.Get("role")
	role, _ := roleInterface.(string)
	userIDInterface, _ := ctx.Get("user_id")
	userID, _ := userIDInterface.(uint)

	// Jika user adalah admin, langsung ambil data
	if role == "admin" {
		response, err := c.useCase.GetAirMonitoringData(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "GetMonitoringAir retrieved successfully",
			"data":    response,
		})
		return
	}

	// Cek apakah user memiliki akses ke setting_id
	hakAksesList, err := c.hakAksesUc.GetHakAksesByGedungIDUser(id, int(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika tidak ditemukan hak akses dengan setting_id dan user_id yang cocok, tolak akses
	if len(hakAksesList) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Jika user memiliki akses, ambil data air monitoring
	response, err := c.useCase.GetAirMonitoringData(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "GetMonitoringAir retrieved successfully",
		"data":    response,
	})
}

func (c *MonitoringDataController) GetListrikMonitoringData(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	roleInterface, _ := ctx.Get("role")
	role, _ := roleInterface.(string)
	userIDInterface, _ := ctx.Get("user_id")
	userID, _ := userIDInterface.(uint)

	// Jika user adalah admin, langsung ambil data
	if role == "admin" {
		response, err := c.useCase.GetListrikMonitoringData(id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "GetMonitoringListrik retrieved successfully",
			"data":    response,
		})
		return
	}

	// Cek apakah user memiliki akses ke setting_id
	hakAksesList, err := c.hakAksesUc.GetHakAksesByGedungIDUser(id, int(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika tidak ditemukan hak akses dengan setting_id dan user_id yang cocok, tolak akses
	if len(hakAksesList) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Jika user memiliki akses, ambil data air monitoring
	response, err := c.useCase.GetListrikMonitoringData(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "GetMonitoringListrik retrieved successfully",
		"data":    response,
	})
}
