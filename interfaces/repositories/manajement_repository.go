package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
)

type ManajementRepository interface {
	FindAll() ([]entities.Manajement, error)
	FindByID(NIP uint) (entities.Manajement, error)
	Create(manajement entities.Manajement) (entities.Manajement, error)
	Update(NIP uint, manajement entities.Manajement) (entities.Manajement, error)
	Delete(NIP uint) error
	WithTransaction() *gorm.DB
}
