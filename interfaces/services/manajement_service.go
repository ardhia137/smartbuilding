package services

import (
	"smartbuilding/entities"
)

type ManajementService interface {
	GetAllManajement() ([]entities.ManajementResponse, error)
	GetManajementByID(NIP uint) (entities.ManajementResponse, error)
	CreateManajement(request entities.CreateManajementRequest) (entities.ManajementResponse, error)
	UpdateManajement(NIP uint, request entities.UpdateManajementRequest) (entities.ManajementResponse, error)
	DeleteManajement(NIP uint) error
}
