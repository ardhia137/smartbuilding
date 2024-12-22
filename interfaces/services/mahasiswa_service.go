package services

import (
	"smartbuilding/entities"
)

type MahasiswaService interface {
	GetAllMahasiswa() ([]entities.MahasiswaResponse, error)
	GetMahasiswaByID(id uint) (entities.MahasiswaResponse, error)
	CreateMahasiswa(request entities.CreateMahasiswaRequest) (entities.MahasiswaResponse, error)
	UpdateMahasiswa(NPM uint, request entities.UpdateMahasiswaRequest) (entities.MahasiswaResponse, error)
	DeleteMahasiswa(NPM uint) error
}
