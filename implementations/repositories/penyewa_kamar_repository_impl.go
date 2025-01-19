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

// NewPenyewaKamarRepository membuat instance baru PenyewaKamarRepository
func NewPenyewaKamarRepository(db *gorm.DB) repositories.PenyewaKamarRepository {
	return &penyewaKamarRepositoryImpl{db}
}

// FindAll mengambil semua penyewaKamar dari database
func (r *penyewaKamarRepositoryImpl) FindAll() ([]entities.PenyewaKamar, error) {
	var penyewaKamars []entities.PenyewaKamar
	err := r.db.Find(&penyewaKamars).Error
	return penyewaKamars, err
}

// CountByKamarID menghitung jumlah penyewa untuk kamar tertentu
//
//	func (r *penyewaKamarRepositoryImpl) CountByKamarID(kamarID uint) (int64, error) {
//		var count int64
//		// Menggunakan GORM untuk menghitung jumlah data dengan kondisi tertentu
//		err := r.db.Model(&entities.PenyewaKamar{}).Where("kamar_id = ?", kamarID).Count(&count).Error
//		if err != nil {
//			return 0, err
//		}
//		return count, nil
//	}
func (r *penyewaKamarRepositoryImpl) CountAktifByKamarID(kamarID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entities.PenyewaKamar{}).Where("kamar_id = ? AND status = ?", kamarID, "tinggal").Count(&count).Error
	fmt.Println("a", count)
	return count, err
}

// FindByID mencari penyewaKamar berdasarkan ID
func (r *penyewaKamarRepositoryImpl) FindByID(id uint) (entities.PenyewaKamar, error) {
	var penyewaKamar entities.PenyewaKamar
	err := r.db.First(&penyewaKamar, id).Error
	return penyewaKamar, err
}

// Create membuat penyewaKamar baru di database
func (r *penyewaKamarRepositoryImpl) Create(penyewaKamar entities.PenyewaKamar) (entities.PenyewaKamar, error) {
	err := r.db.Create(&penyewaKamar).Error
	return penyewaKamar, err
}
func (r *penyewaKamarRepositoryImpl) FindByNPM(npm uint) (entities.PenyewaKamar, error) {
	var penyewaKamar entities.PenyewaKamar
	err := r.db.Where("npm = ? and status = ?", npm, "tinggal").First(&penyewaKamar).Error
	if err != nil {
		return entities.PenyewaKamar{}, err // Kembalikan penyewaKamar kosong dan error jika tidak ditemukan
	}
	return penyewaKamar, nil // Kembalikan penyewaKamar yang ditemukan dan error nil
}

// Update memperbarui informasi penyewaKamar berdasarkan ID
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

// Delete menghapus penyewaKamar dari database berdasarkan ID
func (r *penyewaKamarRepositoryImpl) Delete(id uint) error {
	// Cari penyewaKamar berdasarkan ID
	var penyewaKamar entities.PenyewaKamar
	err := r.db.First(&penyewaKamar, id).Error
	if err != nil {
		return err // Jika penyewaKamar tidak ditemukan
	}

	// Hapus penyewaKamar
	err = r.db.Delete(&penyewaKamar).Error
	return err
}
