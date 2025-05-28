package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
)

type UserRepository interface {
	FindAll(role string, user_id uint) ([]entities.User, error)
	FindByID(id uint) (entities.User, error)
	Create(user entities.User) (entities.User, error)
	Update(id uint, user entities.User) (entities.User, error)
	Delete(id uint) error
	WithTransaction() *gorm.DB
}
