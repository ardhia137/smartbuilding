package usecases

import "smartbuilding/entities"

type TorentUseCase interface {
	CreateTorent(request entities.CreateTorentRequest) (*entities.TorentResponse, error)
	GetAllTorent() ([]entities.TorentResponse, error)
	GetTorentByID(id int) (*entities.TorentResponse, error)
	GetTorentByGedungID(id int) ([]entities.TorentResponse, error)
	UpdateTorent(id int, request entities.CreateTorentRequest) (*entities.TorentResponse, error)
	DeleteTorent(id int) error
}
