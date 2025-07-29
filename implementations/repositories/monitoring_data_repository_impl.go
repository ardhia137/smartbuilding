package repositories

import (
	"smartbuilding/entities"

	"gorm.io/gorm"
)

type MonitoringDataRepositoryImpl struct {
	db *gorm.DB
}

func NewMonitoringDataRepository(db *gorm.DB) *MonitoringDataRepositoryImpl {
	return &MonitoringDataRepositoryImpl{db: db}
}

func (r *MonitoringDataRepositoryImpl) SaveMonitoringData(monitoringData entities.MonitoringData) (entities.MonitoringData, error) {
	if err := r.db.Create(&monitoringData).Error; err != nil {
		return entities.MonitoringData{}, err
	}
	return monitoringData, nil
}

func (r *MonitoringDataRepositoryImpl) GetAirMonitoringData(id int) ([]entities.MonitoringData, error) {
	var monitoringData []entities.MonitoringData
	if err := r.db.Where("monitoring_name LIKE ?", "monitoring_air%").Where("id_gedung = ?", id).Find(&monitoringData).Error; err != nil {
		return nil, err
	}
	return monitoringData, nil
}

func (r *MonitoringDataRepositoryImpl) GetAirMonitoringDataHarian(id int) ([]entities.MonitoringData, error) {
	var monitoringDataHarian []entities.MonitoringData
	if err := r.db.Table("monitoring_data_harian").Where("monitoring_name LIKE ?", "monitoring_air%").Where("id_gedung = ?", id).Find(&monitoringDataHarian).Error; err != nil {
		return nil, err
	}
	return monitoringDataHarian, nil
}

func (r *MonitoringDataRepositoryImpl) GetListrikMonitoringData(id int) ([]entities.MonitoringData, error) {
	var monitoringData []entities.MonitoringData
	if err := r.db.Where("monitoring_name LIKE ?", "monitoring_listrik_%").Where("id_gedung = ?", id).Find(&monitoringData).Error; err != nil {
		return nil, err
	}
	return monitoringData, nil
}

func (r *MonitoringDataRepositoryImpl) GetListrikMonitoringDataHarian(id int) ([]entities.MonitoringData, error) {
	var monitoringDataHarian []entities.MonitoringData
	if err := r.db.Table("monitoring_data_harian").Where("monitoring_name LIKE ?", "monitoring_listrik_%").Where("id_gedung = ?", id).Find(&monitoringDataHarian).Error; err != nil {
		return nil, err
	}
	return monitoringDataHarian, nil
}

func (r *MonitoringDataRepositoryImpl) FindAll() ([]entities.MonitoringData, error) {
	var data []entities.MonitoringData
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MonitoringDataRepositoryImpl) FindBySettingID(id int) ([]entities.MonitoringData, error) {
	var data []entities.MonitoringData
	if err := r.db.Find(&data).Where("id_gedung = ?", id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MonitoringDataRepositoryImpl) SaveHarianData(data entities.MonitoringData) (*entities.MonitoringData, error) {
	if err := r.db.Table("monitoring_data_harian").Create(&data).Error; err != nil {
		return nil, err
	}
	return &data, nil
}

func (r *MonitoringDataRepositoryImpl) Truncate() error {
	if err := r.db.Exec("TRUNCATE TABLE monitoring_data").Error; err != nil {
		return err
	}
	return nil
}
