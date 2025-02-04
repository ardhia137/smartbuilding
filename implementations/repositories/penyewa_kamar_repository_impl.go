package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type penyewaKamarRepositoryImpl struct {
	db *gorm.DB
}

func NewPenyewaKamarRepository(db *gorm.DB) repositories.PenyewaKamarRepository {
	return &penyewaKamarRepositoryImpl{db}
}

func (r *penyewaKamarRepositoryImpl) FindAll() ([]entities.PenyewaKamar, error) {
	var penyewaKamars []entities.PenyewaKamar
	err := r.db.Find(&penyewaKamars).Error
	return penyewaKamars, err
}

func (r *penyewaKamarRepositoryImpl) CountAktifByKamarID(kamarID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.PenyewaKamar{}).Where("kamar_id = ? AND status = ?", kamarID, "tinggal").Count(&count).Error
	fmt.Println("a", count)
	return count, err
}

func (r *penyewaKamarRepositoryImpl) FindByID(id uint) (entities.PenyewaKamar, error) {
	var penyewaKamar entities.PenyewaKamar
	err := r.db.First(&penyewaKamar, id).Error
	return penyewaKamar, err
}

func (r *penyewaKamarRepositoryImpl) Create(penyewaKamar entities.PenyewaKamar) (entities.PenyewaKamar, error) {
	err := r.db.Create(&penyewaKamar).Error
	return penyewaKamar, err
}

func (r *penyewaKamarRepositoryImpl) FindByNPM(npm uint) (entities.PenyewaKamar, error) {
	var penyewaKamar entities.PenyewaKamar
	err := r.db.Where("npm = ? and status = ?", npm, "tinggal").First(&penyewaKamar).Error
	if err != nil {
		return entities.PenyewaKamar{}, err
	}
	return penyewaKamar, nil
}

func (r *penyewaKamarRepositoryImpl) Update(id uint, penyewaKamar entities.PenyewaKamar) (entities.PenyewaKamar, error) {
	var existingPenyewaKamar entities.PenyewaKamar
	err := r.db.First(&existingPenyewaKamar, id).Error
	if err != nil {
		return entities.PenyewaKamar{}, err
	}
	penyewaKamar.ID = existingPenyewaKamar.ID
	err = r.db.Save(&penyewaKamar).Error
	return penyewaKamar, err
}

func (r *penyewaKamarRepositoryImpl) Delete(id uint) error {
	var penyewaKamar entities.PenyewaKamar
	err := r.db.First(&penyewaKamar, id).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(&penyewaKamar).Error
	return err
}
