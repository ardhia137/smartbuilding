package usecases

import "smartbuilding/entities"

type PengelolaGedungUseCase interface {
	CreatePengelolaGedung(request entities.CreatePengelolaGedungRequest) (*entities.PengelolaGedungResponse, error)
	GetAllPengelolaGedung() ([]entities.AllPengelolaGedungResponse, error)
	GetPengelolaGedungByID(id int) (*entities.PengelolaGedungResponse, error)
	GetPengelolaGedungBySettingIDUser(id int, userId int) ([]entities.PengelolaGedungResponse, error)
	GetPengelolaGedungByUser(userId int) ([]entities.AllPengelolaGedungResponse, error)
	UpdatePengelolaGedung(id int, request entities.CreatePengelolaGedungRequest) (*entities.PengelolaGedungResponse, error)
	DeletePengelolaGedung(id int) error
}
