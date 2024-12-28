package controllers

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"net/http"
	"smartbuilding/entities"
	"smartbuilding/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MahasiswaController struct {
	mahasiswaUseCase usecases.MahasiswaUseCase
}

func NewMahasiswaController(mahasiswaUseCase usecases.MahasiswaUseCase) *MahasiswaController {
	return &MahasiswaController{mahasiswaUseCase: mahasiswaUseCase}
}

func (c *MahasiswaController) GetAllMahasiswas(ctx *gin.Context) {
	mahasiswas, err := c.mahasiswaUseCase.GetAllMahasiswa()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to retrieve mahasiswas",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mahasiswas retrieved successfully",
		"data":    mahasiswas,
	})
}

func (c *MahasiswaController) GetMahasiswaByID(ctx *gin.Context) {
	NPM, err := strconv.Atoi(ctx.Param("NPM"))
	fmt.Println("Received NPM parameter:", ctx.Param("NPM"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": NPM,
		})
		return
	}
	mahasiswa, err := c.mahasiswaUseCase.GetMahasiswaByID(uint(NPM))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Mahasiswa not found",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mahasiswa retrieved successfully",
		"data":    mahasiswa,
	})
}
func (c *MahasiswaController) CreateMahasiswa(ctx *gin.Context) {
	var request entities.CreateMahasiswaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	mahasiswa, err := c.mahasiswaUseCase.CreateMahasiswa(request)
	if err != nil {
		// Tangani error spesifik
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			ctx.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "Duplicate entry error",
				"error":   mysqlErr.Message,
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create mahasiswa",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Mahasiswa created successfully",
		"data":    mahasiswa,
	})
}

func (c *MahasiswaController) UpdateMahasiswa(ctx *gin.Context) {
	NPM, err := strconv.Atoi(ctx.Param("NPM"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid mahasiswa ID format",
		})
		return
	}
	var request entities.UpdateMahasiswaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}
	mahasiswa, err := c.mahasiswaUseCase.UpdateMahasiswa(uint(NPM), request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update mahasiswa",
			"error":   err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mahasiswa updated successfully",
		"data":    mahasiswa,
	})
}

func (c *MahasiswaController) DeleteMahasiswa(ctx *gin.Context) {
	NPM, err := strconv.Atoi(ctx.Param("NPM"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid mahasiswa ID format",
		})
		return
	}

	err = c.mahasiswaUseCase.DeleteMahasiswa(uint(NPM))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete mahasiswa",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mahasiswa deleted successfully",
	})
}
