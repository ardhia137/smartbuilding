package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type DataTorenRepositoryImpl struct {
	db *gorm.DB
}

func NewDataTorenRepository(db *gorm.DB) repositories.DataTorenRepository {
	return &DataTorenRepositoryImpl{db: db}
}

func (r *DataTorenRepositoryImpl) Create(dataToren *entities.DataToren) (*entities.DataToren, error) {
	if err := r.db.Create(dataToren).Error; err != nil {
		return nil, err
	}
	return dataToren, nil
}

func (r *DataTorenRepositoryImpl) FindAll() ([]entities.DataToren, error) {
	var dataTorenList []entities.DataToren
	if err := r.db.Find(&dataTorenList).Error; err != nil {
		return nil, err
	}
	return dataTorenList, nil
}

func (r *DataTorenRepositoryImpl) FindBySettingID(id int) ([]entities.DataToren, error) {
	var dataTorenList []entities.DataToren
	if err := r.db.Find(&dataTorenList).Where("id_setting = ?", id).Error; err != nil {
		return nil, err
	}
	return dataTorenList, nil
}

func (r *DataTorenRepositoryImpl) FindByID(id int) (*entities.DataToren, error) {
	var dataToren entities.DataToren
	if err := r.db.First(&dataToren, id).Error; err != nil {
		return nil, err
	}
	return &dataToren, nil
}

func (r *DataTorenRepositoryImpl) Update(dataToren *entities.DataToren) (*entities.DataToren, error) {
	if err := r.db.Save(dataToren).Error; err != nil {
		return nil, err
	}
	return dataToren, nil
}

func (r *DataTorenRepositoryImpl) Delete(id int) error {
	if err := r.db.Delete(&entities.DataToren{}, id).Error; err != nil {
		return err
	}
	return nil
}
