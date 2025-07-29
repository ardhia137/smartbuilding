package controllers

import (
	"net/http"
	"strconv"

	"smartbuilding/entities"
	"smartbuilding/usecases"

	"github.com/gin-gonic/gin"
)

type HakAksesController struct {
	hakAksesUseCase usecases.HakAksesUseCase
}

func NewHakAksesController(hakAksesUseCase usecases.HakAksesUseCase) *HakAksesController {
	return &HakAksesController{hakAksesUseCase: hakAksesUseCase}
}

func (c *HakAksesController) CreateHakAkses(ctx *gin.Context) {
	var request entities.CreateHakAksesRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.hakAksesUseCase.CreateHakAkses(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "HakAkses created successfully",
		"data":    response,
	})
}

func (c *HakAksesController) GetAllHakAkses(ctx *gin.Context) {
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

	var (
		response interface{}
		err      error
	)
	if role == "admin" {
		response, err = c.hakAksesUseCase.GetAllHakAkses()
	} else {
		response, err = c.hakAksesUseCase.GetHakAksesByUser(int(userID))
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data hak akses tidak ditemukan"})
		return
	}

	// Jika response berupa slice kosong
	switch v := response.(type) {
	case []entities.HakAkses:
		if len(v) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Data hak akses tidak ditemukan"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "HakAkses retrieved successfully",
		"data":    response,
	})
}

func (c *HakAksesController) GetHakAksesByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := c.hakAksesUseCase.GetHakAksesByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "HakAkses retrieved successfully",
		"data":    response,
	})
}

func (c *HakAksesController) GetHakAksesByGedungID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
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

	response, err := c.hakAksesUseCase.GetHakAksesByGedungIDUser(id, int(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "HakAkses retrieved successfully",
		"data":    response,
	})
}

func (c *HakAksesController) UpdateHakAkses(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var request entities.CreateHakAksesRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.hakAksesUseCase.UpdateHakAkses(id, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "HakAkses updated successfully",
		"data":    response,
	})
}

func (c *HakAksesController) DeleteHakAkses(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.hakAksesUseCase.DeleteHakAkses(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Hak Akses deleted successfully",
	})
}
