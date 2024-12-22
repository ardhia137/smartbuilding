package repositories

import (
	"smartbuilding/entities"
)

type UserRepository interface {
	FindAll() ([]entities.User, error)
	FindByID(id uint) (entities.User, error)
	Create(user entities.User) (entities.User, error)
	Update(id uint, user entities.User) (entities.User, error)
	Delete(id uint) error
}
