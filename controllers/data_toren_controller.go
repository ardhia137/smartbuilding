package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"smartbuilding/entities"
	"smartbuilding/usecases"
)

type DataTorenController struct {
	dataTorenUseCase usecases.DataTorenUseCase
}

func NewDataTorenController(dataTorenUseCase usecases.DataTorenUseCase) *DataTorenController {
	return &DataTorenController{dataTorenUseCase: dataTorenUseCase}
}

func (c *DataTorenController) CreateDataToren(ctx *gin.Context) {
	var request entities.CreateDataTorenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.dataTorenUseCase.CreateDataToren(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *DataTorenController) GetAllDataToren(ctx *gin.Context) {
	response, err := c.dataTorenUseCase.GetAllDataToren()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DataTorenController) GetDataTorenByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := c.dataTorenUseCase.GetDataTorenByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DataTorenController) GetDataTorenBySettingID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := c.dataTorenUseCase.GetDataTorenBySettingID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DataTorenController) UpdateDataToren(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var request entities.CreateDataTorenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.dataTorenUseCase.UpdateDataToren(id, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *DataTorenController) DeleteDataToren(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.dataTorenUseCase.DeleteDataToren(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
