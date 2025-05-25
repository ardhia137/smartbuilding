package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"smartbuilding/entities"
	"smartbuilding/usecases"
)

type PengelolaGedungController struct {
	pengelolaGedungUseCase usecases.PengelolaGedungUseCase
}

func NewPengelolaGedungController(pengelolaGedungUseCase usecases.PengelolaGedungUseCase) *PengelolaGedungController {
	return &PengelolaGedungController{pengelolaGedungUseCase: pengelolaGedungUseCase}
}

func (c *PengelolaGedungController) CreatePengelolaGedung(ctx *gin.Context) {
	var request entities.CreatePengelolaGedungRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.pengelolaGedungUseCase.CreatePengelolaGedung(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}
func (c *PengelolaGedungController) GetAllPengelolaGedung(ctx *gin.Context) {
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
		response, err = c.pengelolaGedungUseCase.GetAllPengelolaGedung()
	} else {
		response, err = c.pengelolaGedungUseCase.GetPengelolaGedungByUser(int(userID))
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data pengelola gedung tidak ditemukan"})
		return
	}

	// Jika response berupa slice kosong
	switch v := response.(type) {
	case []entities.PengelolaGedung:
		if len(v) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Data pengelola gedung tidak ditemukan"})
			return
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PengelolaGedungController) GetPengelolaGedungByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	response, err := c.pengelolaGedungUseCase.GetPengelolaGedungByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PengelolaGedungController) GetPengelolaGedungBySettingID(ctx *gin.Context) {
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

	var response interface{}

	if role == "admin" {
		response, err = c.pengelolaGedungUseCase.GetPengelolaGedungByID(id)
	} else {
		response, err = c.pengelolaGedungUseCase.GetPengelolaGedungBySettingIDUser(id, int(userID))
	}

	// Cek jika terjadi error
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Cek jika response kosong (nil)
	if response == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Data pengelola gedung tidak ditemukan"})
		return
	}

	// Cek panjang slice jika response berupa slice
	if dataSlice, ok := response.([]entities.PengelolaGedung); ok {
		if len(dataSlice) == 0 {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Data pengelola gedung tidak ditemukan"})
			return
		}
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PengelolaGedungController) UpdatePengelolaGedung(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var request entities.CreatePengelolaGedungRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.pengelolaGedungUseCase.UpdatePengelolaGedung(id, request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PengelolaGedungController) DeletePengelolaGedung(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	if err := c.pengelolaGedungUseCase.DeletePengelolaGedung(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pengelola Gedung deleted successfully",
	})
}
