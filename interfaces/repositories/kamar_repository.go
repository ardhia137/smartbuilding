package repositories

import (
	"smartbuilding/entities"
)

type KamarRepository interface {
	FindAll() ([]entities.Kamar, error)
	FindByID(id uint) (entities.Kamar, error)
	Create(kamar entities.Kamar) (entities.Kamar, error)
	Update(id uint, kamar entities.Kamar) (entities.Kamar, error)
	Delete(id uint) error
}
