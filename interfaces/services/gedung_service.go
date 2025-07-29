package services

import (
	"smartbuilding/entities"
)

type GedungService interface {
	CreateGedung(request entities.CreateGedungRequest) (*entities.GedungResponseCreate, error)
	GetAllGedung(role string, userID uint) ([]entities.GedungResponse, error)
	GetAllCornJobs() ([]entities.GedungResponse, error)
	GetGedungByID(id int) (*entities.GedungResponse, error)
	UpdateGedung(id int, request entities.CreateGedungRequest) (*entities.GedungResponse, error)
	DeleteGedung(id int) error
}
