package services

import (
	"smartbuilding/entities"
)

type KamarService interface {
	GetAllKamar() ([]entities.KamarResponse, error)
	GetKamarByID(id uint) (entities.KamarResponse, error)
	CreateKamar(request entities.CreateKamarRequest) (entities.KamarResponse, error)
	UpdateKamar(id uint, request entities.CreateKamarRequest) (entities.KamarResponse, error)
	DeleteKamar(id uint) error
}