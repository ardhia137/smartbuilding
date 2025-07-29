package repositories

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"

	"gorm.io/gorm"
)

type GedungRepositoryImpl struct {
	db *gorm.DB
}

func (r *GedungRepositoryImpl) WithTransaction() *gorm.DB {
	return r.db
}

func NewGedungRepository(db *gorm.DB) repositories.GedungRepository {
	return &GedungRepositoryImpl{db: db}
}

func (r *GedungRepositoryImpl) Create(gedung *entities.Gedung) (*entities.Gedung, error) {
	if err := r.db.Create(gedung).Error; err != nil {
		return nil, err
	}
	return gedung, nil
}

func (r *GedungRepositoryImpl) FindAll() ([]entities.Gedung, error) {
	var gedungList []entities.Gedung
	if err := r.db.Find(&gedungList).Error; err != nil {
		return nil, err
	}
	return gedungList, nil
}

func (r *GedungRepositoryImpl) FindByID(id int) (*entities.Gedung, error) {
	var gedung entities.Gedung
	if err := r.db.First(&gedung, id).Error; err != nil {
		return nil, err
	}
	return &gedung, nil
}

func (r *GedungRepositoryImpl) FindByUserId(userID uint) ([]entities.Gedung, error) {
	var gedungList []entities.Gedung

	err := r.db.
		Joins("JOIN hak_akses ha ON setting.id = ha.gedung_id").
		Joins("JOIN user u ON u.id = ha.user_id").
		Where("u.id = ?", userID).
		Find(&gedungList).Error

	if err != nil {
		return nil, err
	}

	return gedungList, nil
}

func (r *GedungRepositoryImpl) Update(gedung *entities.Gedung) (*entities.Gedung, error) {
	if err := r.db.Save(gedung).Error; err != nil {
		return nil, err
	}
	return gedung, nil
}

func (r *GedungRepositoryImpl) Delete(id int) error {
	if err := r.db.Delete(&entities.Gedung{}, id).Error; err != nil {
		return err
	}
	return nil
}
