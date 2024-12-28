package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type manajementRepositoryImpl struct {
	db *gorm.DB
}

// NewManajementRepository membuat instance baru ManajementRepository
func NewManajementRepository(db *gorm.DB) repositories.ManajementRepository {
	return &manajementRepositoryImpl{db}
}

// FindAll mengambil semua manajement dari database
func (r *manajementRepositoryImpl) FindAll() ([]entities.Manajement, error) {
	var manajements []entities.Manajement
	err := r.db.Preload("User").Find(&manajements).Error
	return manajements, err
}

// FindByID mencari manajement berdasarkan ID
func (r *manajementRepositoryImpl) FindByID(NIP uint) (entities.Manajement, error) {
	var manajement entities.Manajement
	err := r.db.Preload("User").First(&manajement, NIP).Error
	return manajement, err
}

// Create membuat manajement baru di database
func (r *manajementRepositoryImpl) Create(manajement entities.Manajement) (entities.Manajement, error) {
	err := r.db.Create(&manajement).Error
	return manajement, err
}
func (r *manajementRepositoryImpl) WithTransaction() *gorm.DB {
	return r.db
}

// Update memperbarui informasi manajement berdasarkan ID
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

// Delete menghapus manajement dari database berdasarkan ID
func (r *manajementRepositoryImpl) Delete(NIP uint) error {
	// Cari manajement berdasarkan ID
	var manajement entities.Manajement
	err := r.db.First(&manajement, NIP).Error
	if err != nil {
		return err // Jika manajement tidak ditemukan
	}
	// Hapus manajement
	err = r.db.Delete(&manajement).Error
	return err
}
