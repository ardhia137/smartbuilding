package repositories

import (
	"smartbuilding/entities"
)

type PenyewaKamarRepository interface {
	FindAll() ([]entities.PenyewaKamar, error)
	FindByID(id uint) (entities.PenyewaKamar, error)
	CountAktifByKamarID(kamarID uint) (int64, error)
	FindByNPM(npm uint) (entities.PenyewaKamar, error)
	Create(penyewaKamar entities.PenyewaKamar) (entities.PenyewaKamar, error)
	Update(id uint, penyewaKamar entities.PenyewaKamar) (entities.PenyewaKamar, error)
	Delete(id uint) error
}
