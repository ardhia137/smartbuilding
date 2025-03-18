package repositories

import (
	"smartbuilding/entities"
)

type PengelolaGedungRepository interface {
	Create(dataToren *entities.PengelolaGedung) (*entities.PengelolaGedung, error)
	FindAll() ([]entities.AllPengelolaGedungResponse, error)
	FindByID(id int) (*entities.PengelolaGedung, error)
	FindBySettingIDUser(id int, userID int) ([]entities.PengelolaGedung, error)
	FindBySettingUser(userID int) ([]entities.AllPengelolaGedungResponse, error)
	Update(dataToren *entities.PengelolaGedung) (*entities.PengelolaGedung, error)
	Delete(id int) error
}
