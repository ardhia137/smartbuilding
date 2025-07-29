package services

import "smartbuilding/entities"

type HakAksesService interface {
	CreateHakAkses(request entities.CreateHakAksesRequest) (*entities.HakAksesResponse, error)
	GetAllHakAkses() ([]entities.AllHakAksesResponse, error)
	GetHakAksesByID(id int) (*entities.HakAksesResponse, error)
	GetHakAksesByGedungIDUser(id int, userId int) ([]entities.HakAksesResponse, error)
	GetHakAksesByUser(userId int) ([]entities.AllHakAksesResponse, error)
	UpdateHakAkses(id int, request entities.CreateHakAksesRequest) (*entities.HakAksesResponse, error)
	DeleteHakAkses(id int) error
}
