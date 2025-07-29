package repositories

import (
	"smartbuilding/entities"
)

type TorentRepository interface {
	Create(torent *entities.Torent) (*entities.Torent, error)
	FindAll() ([]entities.Torent, error)
	FindByID(id int) (*entities.Torent, error)
	FindByGedungID(id int) ([]entities.Torent, error)
	Update(torent *entities.Torent) (*entities.Torent, error)
	Delete(id int) error
}
