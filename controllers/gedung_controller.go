package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"smartbuilding/entities"
	"smartbuilding/usecases"

	"github.com/gin-gonic/gin"
)

type GedungController struct {
	gedungUseCase usecases.GedungUseCase
}

func NewGedungController(gedungUseCase usecases.GedungUseCase) *GedungController {
	return &GedungController{gedungUseCase: gedungUseCase}
}

func (c *GedungController) CreateGedung(ctx *gin.Context) {
	var request entities.CreateGedungRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.gedungUseCase.CreateGedung(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Gedung created successfully",
		"data":    response,
	})
}

func (c *GedungController) GetAllGedung(ctx *gin.Context) {
	// Ambil role dari context
	roleInterface, _ := ctx.Get("role")
	role, _ := roleInterface.(string)

	// Ambil user_id dari context
	userIDInterface, _ := ctx.Get("user_id")
	userID, _ := userIDInterface.(uint)
	fmt.Println(userID)
	response, err := c.gedungUseCase.GetAllGedung(role, userID)

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
		"message": "Gedung retrieved successfully",
		"data":    response,
	})
}

func (c *GedungController) GetGedungByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := c.gedungUseCase.GetGedungByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Gedung retrieved successfully",
		"data":    response,
	})
}

func (c *GedungController) UpdateGedung(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var request entities.CreateGedungRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.gedungUseCase.UpdateGedung(id, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Gedung updated successfully",
		"data":    response,
	})
}

func (c *GedungController) DeleteGedung(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.gedungUseCase.DeleteGedung(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Gedung deleted successfully",
	})
}
