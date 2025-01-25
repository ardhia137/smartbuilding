package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type manajementRepositoryImpl struct {
	db *gorm.DB
}

func NewManajementRepository(db *gorm.DB) repositories.ManajementRepository {
	return &manajementRepositoryImpl{db}
}

func (r *manajementRepositoryImpl) FindAll() ([]entities.Manajement, error) {
	var manajements []entities.Manajement
	err := r.db.Preload("User").Find(&manajements).Error
	return manajements, err
}

func (r *manajementRepositoryImpl) FindByID(NIP uint) (entities.Manajement, error) {
	var manajement entities.Manajement
	err := r.db.Preload("User").First(&manajement, NIP).Error
	return manajement, err
}

func (r *manajementRepositoryImpl) Create(manajement entities.Manajement) (entities.Manajement, error) {
	err := r.db.Create(&manajement).Error
	return manajement, err
}

func (r *manajementRepositoryImpl) WithTransaction() *gorm.DB {
	return r.db
}

func (r *manajementRepositoryImpl) Update(NIP uint, manajement entities.Manajement) (entities.Manajement, error) {
	var existingManajement entities.Manajement
	err := r.db.First(&existingManajement, NIP).Error
	if err != nil {
		return entities.Manajement{}, err
	}
	manajement.NIP = existingManajement.NIP
	err = r.db.Save(&manajement).Error
	return manajement, err
}

func (r *manajementRepositoryImpl) Delete(NIP uint) error {
	var manajement entities.Manajement
	err := r.db.First(&manajement, NIP).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(&manajement).Error
	return err
}
