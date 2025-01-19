package controllers

import (
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PenyewaKamarController struct {
	penyewaPenyewaKamarUseCase usecases.PenyewaKamarUseCase
}

func NewPenyewaKamarController(penyewaPenyewaKamarUseCase usecases.PenyewaKamarUseCase) *PenyewaKamarController {
	return &PenyewaKamarController{penyewaPenyewaKamarUseCase: penyewaPenyewaKamarUseCase}
}

func (c *PenyewaKamarController) GetAllPenyewaKamars(ctx *gin.Context) {
	penyewaPenyewaKamars, err := c.penyewaPenyewaKamarUseCase.GetAllPenyewaKamar()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve penyewaPenyewaKamars",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "PenyewaKamars retrieved successfully",
		"data":    penyewaPenyewaKamars,
	})
}

func (c *PenyewaKamarController) GetPenyewaKamarByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid penyewaPenyewaKamar ID format",
		})
		return
	}
	penyewaPenyewaKamar, err := c.penyewaPenyewaKamarUseCase.GetPenyewaKamarByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "PenyewaKamar not found",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "PenyewaKamar retrieved successfully",
		"data":    penyewaPenyewaKamar,
	})
}

func (c *PenyewaKamarController) CreatePenyewaKamar(ctx *gin.Context) {
	var request entities.CreatePenyewaKamarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	penyewaPenyewaKamar, err := c.penyewaPenyewaKamarUseCase.CreatePenyewaKamar(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create penyewaPenyewaKamar",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "PenyewaKamar created successfully",
		"data":    penyewaPenyewaKamar,
	})
}

func (c *PenyewaKamarController) UpdatePenyewaKamar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid penyewaPenyewaKamar ID format",
		})
		return
	}
	var request entities.UpdatePenyewaKamarRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	penyewaPenyewaKamar, err := c.penyewaPenyewaKamarUseCase.UpdatePenyewaKamar(uint(id), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update penyewaPenyewaKamar",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "PenyewaKamar updated successfully",
		"data":    penyewaPenyewaKamar,
	})
}

func (c *PenyewaKamarController) DeletePenyewaKamar(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid penyewaPenyewaKamar ID format",
		})
		return
	}

	err = c.penyewaPenyewaKamarUseCase.DeletePenyewaKamar(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete penyewaPenyewaKamar",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "PenyewaKamar deleted successfully",
	})
}
