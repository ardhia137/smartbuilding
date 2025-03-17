package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"
)

type MonitoringDataController struct {
	useCase usecases.MonitoringDataUseCase
}

func NewMonitoringDataController(useCase usecases.MonitoringDataUseCase) *MonitoringDataController {
	return &MonitoringDataController{useCase}
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
	response, err := c.useCase.GetAirMonitoringData(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *MonitoringDataController) GetListrikMonitoringData(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	response, err := c.useCase.GetListrikMonitoringData(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
