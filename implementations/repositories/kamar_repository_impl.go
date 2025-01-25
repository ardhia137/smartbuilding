package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type kamarRepositoryImpl struct {
	db *gorm.DB
}

func NewKamarRepository(db *gorm.DB) repositories.KamarRepository {
	return &kamarRepositoryImpl{db}
}

func (r *kamarRepositoryImpl) FindAll() ([]entities.Kamar, error) {
	var kamars []entities.Kamar
	err := r.db.Find(&kamars).Error
	return kamars, err
}

func (r *kamarRepositoryImpl) FindByID(id uint) (entities.Kamar, error) {
	var kamar entities.Kamar
	err := r.db.First(&kamar, id).Error
	return kamar, err
}

func (r *kamarRepositoryImpl) Create(kamar entities.Kamar) (entities.Kamar, error) {
	err := r.db.Create(&kamar).Error
	return kamar, err
}

func (r *kamarRepositoryImpl) Update(id uint, kamar entities.Kamar) (entities.Kamar, error) {
	var existingKamar entities.Kamar
	err := r.db.First(&existingKamar, id).Error
	if err != nil {
		return entities.Kamar{}, err
	}
	kamar.ID = existingKamar.ID
	err = r.db.Save(&kamar).Error
	return kamar, err
}

func (r *kamarRepositoryImpl) Delete(id uint) error {
	var kamar entities.Kamar
	err := r.db.First(&kamar, id).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(&kamar).Error
	return err
}
