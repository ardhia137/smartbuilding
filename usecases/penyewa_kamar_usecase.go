package usecases

import "smartbuilding/entities"

type PenyewaKamarUseCase interface {
	GetAllPenyewaKamar() ([]entities.PenyewaKamarResponse, error)
	GetPenyewaKamarByID(id uint) (entities.PenyewaKamarResponse, error)
	FindByNPM(NPM uint) (entities.PenyewaKamarResponse, error)
	CreatePenyewaKamar(request entities.CreatePenyewaKamarRequest) (entities.PenyewaKamarResponse, error)
	UpdatePenyewaKamar(id uint, request entities.UpdatePenyewaKamarRequest) (entities.PenyewaKamarResponse, error)
	DeletePenyewaKamar(id uint) error
}
