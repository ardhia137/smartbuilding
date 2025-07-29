package repositories

import (
	"smartbuilding/entities"
)

type HakAksesRepository interface {
	Create(hakAkses *entities.HakAkses) (*entities.HakAkses, error)
	FindAll() ([]entities.AllHakAksesResponse, error)
	FindByID(id int) (*entities.HakAkses, error)
	FindByGedungIDUser(id int, userID int) ([]entities.HakAkses, error)
	FindByGedungUser(userID int) ([]entities.AllHakAksesResponse, error)
	FindByUser(userID int) ([]entities.HakAkses, error)
	Update(hakAkses *entities.HakAkses) (*entities.HakAkses, error)
	Delete(id int) error
	WithTransaction() interface{}
}
