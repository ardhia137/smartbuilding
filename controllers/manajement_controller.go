package controllers

import (
	"fmt"
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ManajementController struct {
	manajementUseCase usecases.ManajementUseCase
}

func NewManajementController(manajementUseCase usecases.ManajementUseCase) *ManajementController {
	return &ManajementController{manajementUseCase: manajementUseCase}
}

func (c *ManajementController) GetAllManajements(ctx *gin.Context) {
	manajements, err := c.manajementUseCase.GetAllManajement()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve manajements",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Manajements retrieved successfully",
		"data":    manajements,
	})
}

func (c *ManajementController) GetManajementByID(ctx *gin.Context) {
	NIP, err := strconv.Atoi(ctx.Param("NIP"))
	fmt.Println("Received NIP parameter:", ctx.Param("NIP"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": NIP,
		})
		return
	}
	manajement, err := c.manajementUseCase.GetManajementByID(uint(NIP))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Manajement not found",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Manajement retrieved successfully",
		"data":    manajement,
	})
}

func (c *ManajementController) CreateManajement(ctx *gin.Context) {
	var request entities.CreateManajementRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	manajement, err := c.manajementUseCase.CreateManajement(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create manajement",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Manajement created successfully",
		"data":    manajement,
	})
}

func (c *ManajementController) UpdateManajement(ctx *gin.Context) {
	NIP, err := strconv.Atoi(ctx.Param("NIP"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid manajement ID format",
		})
		return
	}
	var request entities.UpdateManajementRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	manajement, err := c.manajementUseCase.UpdateManajement(uint(NIP), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update manajement",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Manajement updated successfully",
		"data":    manajement,
	})
}

func (c *ManajementController) DeleteManajement(ctx *gin.Context) {
	NIP, err := strconv.Atoi(ctx.Param("NIP"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid manajement ID format",
		})
		return
	}

	err = c.manajementUseCase.DeleteManajement(uint(NIP))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete manajement",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Manajement deleted successfully",
	})
}
