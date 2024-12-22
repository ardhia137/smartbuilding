package controllers

import (
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type KamarController struct {
	kamarUseCase usecases.KamarUseCase
}

func NewKamarController(kamarUseCase usecases.KamarUseCase) *KamarController {
	return &KamarController{kamarUseCase: kamarUseCase}
}

func (c *KamarController) GetAllKamars(ctx *gin.Context) {
	kamars, err := c.kamarUseCase.GetAllKamar()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve kamars",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kamars retrieved successfully",
		"data":    kamars,
	})
}

func (c *KamarController) GetKamarByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid kamar ID format",
		})
		return
	}
	kamar, err := c.kamarUseCase.GetKamarByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Kamar not found",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kamar retrieved successfully",
		"data":    kamar,
	})
}

func (c *KamarController) CreateKamar(ctx *gin.Context) {
	var request entities.CreateKamarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	kamar, err := c.kamarUseCase.CreateKamar(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create kamar",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Kamar created successfully",
		"data":    kamar,
	})
}

func (c *KamarController) UpdateKamar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid kamar ID format",
		})
		return
	}
	var request entities.CreateKamarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	kamar, err := c.kamarUseCase.UpdateKamar(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update kamar",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kamar updated successfully",
		"data":    kamar,
	})
}

func (c *KamarController) DeleteKamar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid kamar ID format",
		})
		return
	}

	err = c.kamarUseCase.DeleteKamar(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete kamar",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Kamar deleted successfully",
	})
}
