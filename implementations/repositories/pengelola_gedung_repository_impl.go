package repositories

import (
	"fmt"
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type PengelolaGedungRepositoryImpl struct {
	db *gorm.DB
}

func NewPengelolaGedungRepository(db *gorm.DB) repositories.PengelolaGedungRepository {
	return &PengelolaGedungRepositoryImpl{db: db}
}

func (r *PengelolaGedungRepositoryImpl) Create(pengelolaGedung *entities.PengelolaGedung) (*entities.PengelolaGedung, error) {
	if err := r.db.Create(pengelolaGedung).Error; err != nil {
		return nil, err
	}
	return pengelolaGedung, nil
}

func (r *PengelolaGedungRepositoryImpl) FindAll() ([]entities.AllPengelolaGedungResponse, error) {
	var pengelolaGedungList []entities.AllPengelolaGedungResponse
	err := r.db.Table("pengelola_gedung pg").
		Select("pg.id,s.nama_gedung, u.username, u.email,u.role,pg.setting_id").
		Joins("JOIN setting s ON pg.setting_id = s.id").
		Joins("JOIN user u ON pg.user_id = u.id").
		Scan(&pengelolaGedungList).Error

	if err != nil {
		return nil, err
	}

	return pengelolaGedungList, nil
}

func (r *PengelolaGedungRepositoryImpl) FindBySettingIDUser(id int, userID int) ([]entities.PengelolaGedung, error) {
	var pengelolaGedungList []entities.PengelolaGedung
	err := r.db.
		Where("setting_id = ? AND user_id = ?", id, userID).
		Find(&pengelolaGedungList).
		Error
	if err != nil {
		return nil, err
	}

	return pengelolaGedungList, nil
}

func (r *PengelolaGedungRepositoryImpl) FindByUser(userID int) ([]entities.PengelolaGedung, error) {
	var pengelolaGedungList []entities.PengelolaGedung
	err := r.db.
		Where("user_id = ?", userID).
		Find(&pengelolaGedungList).
		Error
	if err != nil {
		return nil, err
	}
	fmt.Println(pengelolaGedungList)

	return pengelolaGedungList, nil
}

func (r *PengelolaGedungRepositoryImpl) FindBySettingUser(userID int) ([]entities.AllPengelolaGedungResponse, error) {
	var pengelolaGedungList []entities.AllPengelolaGedungResponse
	subQuery := r.db.Table("pengelola_gedung").
		Select("setting_id").
		Where("user_id = ?", userID)

	err := r.db.Table("pengelola_gedung pg").
		Select("pg.id,s.nama_gedung, u.username, u.email,u.role,pg.setting_id").
		Joins("JOIN setting s ON pg.setting_id = s.id").
		Joins("JOIN user u ON pg.user_id = u.id").
		Where("pg.setting_id IN (?)", subQuery).
		Scan(&pengelolaGedungList).Error
	if err != nil {
		return nil, err
	}

	return pengelolaGedungList, nil
}

func (r *PengelolaGedungRepositoryImpl) FindByID(id int) (*entities.PengelolaGedung, error) {
	var pengelolaGedung entities.PengelolaGedung
	if err := r.db.First(&pengelolaGedung, id).Error; err != nil {
		return nil, err
	}
	return &pengelolaGedung, nil
}

func (r *PengelolaGedungRepositoryImpl) Update(pengelolaGedung *entities.PengelolaGedung) (*entities.PengelolaGedung, error) {
	if err := r.db.Save(pengelolaGedung).Error; err != nil {
		return nil, err
	}
	return pengelolaGedung, nil
}

func (r *PengelolaGedungRepositoryImpl) Delete(id int) error {
	if err := r.db.Delete(&entities.PengelolaGedung{}, id).Error; err != nil {
		return err
	}
	return nil
}
