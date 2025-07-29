package repositories

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"

	"gorm.io/gorm"
)

type SettingRepositoryImpl struct {
	db *gorm.DB
}

func (r *SettingRepositoryImpl) WithTransaction() *gorm.DB {
	return r.db
}

func NewSettingRepository(db *gorm.DB) repositories.SettingRepository {
	return &SettingRepositoryImpl{db: db}
}

func (r *SettingRepositoryImpl) Create(setting *entities.Setting) (*entities.Setting, error) {
	if err := r.db.Create(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

func (r *SettingRepositoryImpl) FindAll() ([]entities.Setting, error) {
	var settingList []entities.Setting
	if err := r.db.Find(&settingList).Error; err != nil {
		return nil, err
	}
	return settingList, nil
}

func (r *SettingRepositoryImpl) FindByID(id int) (*entities.Setting, error) {
	var setting entities.Setting
	if err := r.db.First(&setting, id).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *SettingRepositoryImpl) FindByUserId(userID uint) ([]entities.Setting, error) {
	var settingList []entities.Setting

	err := r.db.
		Joins("JOIN hak_akses ha ON setting.id = ha.setting_id").
		Joins("JOIN user u ON u.id = ha.user_id").
		Where("u.id = ?", userID).
		Find(&settingList).Error

	if err != nil {
		return nil, err
	}

	return settingList, nil
}

func (r *SettingRepositoryImpl) Update(setting *entities.Setting) (*entities.Setting, error) {
	if err := r.db.Save(setting).Error; err != nil {
		return nil, err
	}
	return setting, nil
}

func (r *SettingRepositoryImpl) Delete(id int) error {
	if err := r.db.Delete(&entities.Setting{}, id).Error; err != nil {
		return err
	}
	return nil
}
