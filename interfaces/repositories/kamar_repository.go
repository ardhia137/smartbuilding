package repositories

import (
	"smartbuilding/entities"
)

type KamarRepository interface {
	FindAll() ([]entities.Kamar, error)
	FindByID(id uint) (entities.Kamar, error)
	Create(user entities.Kamar) (entities.Kamar, error)
	Update(id uint, user entities.Kamar) (entities.Kamar, error)
	Delete(id uint) error
}
