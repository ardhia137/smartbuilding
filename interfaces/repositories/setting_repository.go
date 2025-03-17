package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
)

type SettingRepository interface {
	Create(haos *entities.Setting) (*entities.Setting, error)
	FindAll() ([]entities.Setting, error)
	FindByID(id int) (*entities.Setting, error)
	Update(haos *entities.Setting) (*entities.Setting, error)
	Delete(id int) error
	WithTransaction() *gorm.DB
}
