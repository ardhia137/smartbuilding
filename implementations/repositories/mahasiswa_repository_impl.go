package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type mahasiswaRepositoryImpl struct {
	db *gorm.DB
}

func NewMahasiswaRepository(db *gorm.DB) repositories.MahasiswaRepository {
	return &mahasiswaRepositoryImpl{db}
}

func (r *mahasiswaRepositoryImpl) FindAll() ([]entities.Mahasiswa, error) {
	var mahasiswas []entities.Mahasiswa
	err := r.db.Preload("User").Find(&mahasiswas).Error
	return mahasiswas, err
}

func (r *mahasiswaRepositoryImpl) FindByID(NPM uint) (entities.Mahasiswa, error) {
	var mahasiswa entities.Mahasiswa
	err := r.db.Preload("User").First(&mahasiswa, NPM).Error
	return mahasiswa, err
}

func (r *mahasiswaRepositoryImpl) Create(mahasiswa entities.Mahasiswa) (entities.Mahasiswa, error) {
	err := r.db.Create(&mahasiswa).Error
	return mahasiswa, err
}

func (r *mahasiswaRepositoryImpl) WithTransaction() *gorm.DB {
	return r.db
}

func (r *mahasiswaRepositoryImpl) Update(NPM uint, mahasiswa entities.Mahasiswa) (entities.Mahasiswa, error) {
	var existingMahasiswa entities.Mahasiswa
	err := r.db.First(&existingMahasiswa, NPM).Error
	if err != nil {
		return entities.Mahasiswa{}, err
	}
	mahasiswa.NPM = existingMahasiswa.NPM
	err = r.db.Save(&mahasiswa).Error
	return mahasiswa, err
}

func (r *mahasiswaRepositoryImpl) Delete(NPM uint) error {
	var mahasiswa entities.Mahasiswa
	err := r.db.First(&mahasiswa, NPM).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(&mahasiswa).Error
	return err
}
