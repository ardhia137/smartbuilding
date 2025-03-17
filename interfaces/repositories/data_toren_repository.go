package repositories

import (
	"smartbuilding/entities"
)

type DataTorenRepository interface {
	Create(dataToren *entities.DataToren) (*entities.DataToren, error)
	FindAll() ([]entities.DataToren, error)
	FindByID(id int) (*entities.DataToren, error)
	FindBySettingID(id int) ([]entities.DataToren, error)
	Update(dataToren *entities.DataToren) (*entities.DataToren, error)
	Delete(id int) error
}
