package repositories

import (
	"errors"
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"

	"gorm.io/gorm"
)

type HakAksesRepositoryImpl struct {
	db *gorm.DB
}

func NewHakAksesRepository(db *gorm.DB) repositories.HakAksesRepository {
	return &HakAksesRepositoryImpl{db: db}
}

func (r *HakAksesRepositoryImpl) Create(hakAkses *entities.HakAkses) (*entities.HakAkses, error) {
	var existing entities.HakAkses
	err := r.db.Where("gedung_id = ? AND user_id = ?",
		hakAkses.GedungID,
		hakAkses.UserId).
		First(&existing).Error

	if err == nil {
		return nil, fmt.Errorf("hak akses untuk gedung ini sudah ada")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := r.db.Create(hakAkses).Error; err != nil {
		return nil, err
	}

	return hakAkses, nil
}

func (r *HakAksesRepositoryImpl) FindAll() ([]entities.AllHakAksesResponse, error) {
	var hakAksesList []entities.AllHakAksesResponse
	err := r.db.Table("hak_akses ha").
		Select("ha.id,g.nama_gedung, u.username, u.email,u.role,ha.gedung_id").
		Joins("JOIN gedung g ON ha.gedung_id = g.id").
		Joins("JOIN user u ON ha.user_id = u.id").
		Scan(&hakAksesList).Error

	if err != nil {
		return nil, err
	}

	return hakAksesList, nil
}

func (r *HakAksesRepositoryImpl) FindByGedungIDUser(id int, userID int) ([]entities.HakAkses, error) {
	var hakAksesList []entities.HakAkses
	if err := r.db.Where("gedung_id = ? AND user_id = ?", id, userID).Find(&hakAksesList).Error; err != nil {
		return nil, err
	}
	return hakAksesList, nil
}

func (r *HakAksesRepositoryImpl) FindByUser(userID int) ([]entities.HakAkses, error) {
	var hakAksesList []entities.HakAkses
	if err := r.db.Where("user_id = ?", userID).Find(&hakAksesList).Error; err != nil {
		return nil, err
	}
	return hakAksesList, nil
}

func (r *HakAksesRepositoryImpl) FindByGedungUser(userID int) ([]entities.AllHakAksesResponse, error) {
	var hakAksesList []entities.AllHakAksesResponse

	subQuery := r.db.Table("hak_akses").
		Select("gedung_id").
		Where("user_id = ?", userID)

	err := r.db.Table("hak_akses ha").
		Select("ha.id,g.nama_gedung, u.username, u.email,u.role,ha.gedung_id").
		Joins("JOIN gedung g ON ha.gedung_id = g.id").
		Joins("JOIN user u ON ha.user_id = u.id").
		Where("ha.gedung_id IN (?)", subQuery).
		Scan(&hakAksesList).Error

	if err != nil {
		return nil, err
	}

	return hakAksesList, nil
}

func (r *HakAksesRepositoryImpl) FindByID(id int) (*entities.HakAkses, error) {
	var hakAkses entities.HakAkses
	if err := r.db.First(&hakAkses, id).Error; err != nil {
		return nil, err
	}
	return &hakAkses, nil
}

func (r *HakAksesRepositoryImpl) Update(hakAkses *entities.HakAkses) (*entities.HakAkses, error) {
	var existing entities.HakAkses
	err := r.db.Where("gedung_id = ? AND user_id = ?",
		hakAkses.GedungID,
		hakAkses.UserId).
		First(&existing).Error

	// Jika record sudah ada dan bukan record yang sama (ID berbeda)
	if err == nil && existing.ID != hakAkses.ID {
		return nil, fmt.Errorf("Hak akses untuk gedung ini sudah ada untuk user: %d, gedung: %d",
			hakAkses.UserId,
			hakAkses.GedungID)
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := r.db.Save(hakAkses).Error; err != nil {
		return nil, err
	}

	return hakAkses, nil
}

func (r *HakAksesRepositoryImpl) Delete(id int) error {
	if err := r.db.Delete(&entities.HakAkses{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *HakAksesRepositoryImpl) WithTransaction() interface{} {
	return r.db
}
