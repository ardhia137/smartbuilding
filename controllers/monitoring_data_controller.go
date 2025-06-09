package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"
)

type MonitoringDataController struct {
	useCase     usecases.MonitoringDataUseCase
	pengelolauc usecases.PengelolaGedungUseCase
}

func NewMonitoringDataController(useCase usecases.MonitoringDataUseCase, pengelolauc usecases.PengelolaGedungUseCase) *MonitoringDataController {
	return &MonitoringDataController{useCase, pengelolauc}
}

// @Summary Menyimpan data monitoring
// @Description Menyimpan data monitoring baru
// @Tags monitoring
// @Accept json
// @Produce json
// @Param data body entities.CreateMonitoringDataRequest true "Data monitoring yang akan disimpan"
// @Success 201 {object} entities.MonitoringDataResponse "Data monitoring berhasil disimpan"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /monitoring-data [post]
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

// @Summary Mendapatkan data monitoring air
// @Description Mendapatkan data monitoring air berdasarkan ID setting
// @Tags monitoring
// @Accept json
// @Produce json
// @Param id path int true "Setting ID"
// @Security ApiKeyAuth
// @Success 200 {array} entities.GetAirDataResponse "Data monitoring air berhasil diambil"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /monitoring-data/air/{id} [get]
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
	pengelolaGedungList, err := c.pengelolauc.GetPengelolaGedungBySettingIDUser(id, int(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika tidak ditemukan pengelola gedung dengan setting_id dan user_id yang cocok, tolak akses
	if len(pengelolaGedungList) == 0 {
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

// @Summary Mendapatkan data monitoring listrik
// @Description Mendapatkan data monitoring listrik berdasarkan ID setting
// @Tags monitoring
// @Accept json
// @Produce json
// @Param id path int true "Setting ID"
// @Security ApiKeyAuth
// @Success 200 {object} entities.GetListrikDataResponse "Data monitoring listrik berhasil diambil"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal Server Error"
// @Router /monitoring-data/listrik/{id} [get]
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
	pengelolaGedungList, err := c.pengelolauc.GetPengelolaGedungBySettingIDUser(id, int(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Jika tidak ditemukan pengelola gedung dengan setting_id dan user_id yang cocok, tolak akses
	if len(pengelolaGedungList) == 0 {
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
