package repositories

import (
	"smartbuilding/entities"

	"gorm.io/gorm"
)

type GedungRepository interface {
	Create(gedung *entities.Gedung) (*entities.Gedung, error)
	FindAll() ([]entities.Gedung, error)
	FindByID(id int) (*entities.Gedung, error)
	FindByUserId(id uint) ([]entities.Gedung, error)
	Update(gedung *entities.Gedung) (*entities.Gedung, error)
	Delete(id int) error
	WithTransaction() *gorm.DB
}
