package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
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

func (r *MonitoringDataRepositoryImpl) GetAirMonitoringData() ([]entities.MonitoringData, error) {
	var monitoringData []entities.MonitoringData
	if err := r.db.Where("monitoring_name LIKE ?", "monitoring_air%").Find(&monitoringData).Error; err != nil {
		return nil, err
	}
	return monitoringData, nil
}

func (r *MonitoringDataRepositoryImpl) GetAirMonitoringDataHarian() ([]entities.MonitoringData, error) {
	var monitoringDataHarian []entities.MonitoringData
	if err := r.db.Table("monitoring_data_harian").Where("monitoring_name LIKE ?", "monitoring_air%").Find(&monitoringDataHarian).Error; err != nil {
		return nil, err
	}
	return monitoringDataHarian, nil
}

func (r *MonitoringDataRepositoryImpl) GetListrikMonitoringData() ([]entities.MonitoringData, error) {
	var monitoringData []entities.MonitoringData
	if err := r.db.Where("monitoring_name LIKE ?", "monitoring_listrik_%").Find(&monitoringData).Error; err != nil {
		return nil, err
	}
	return monitoringData, nil
}

func (r *MonitoringDataRepositoryImpl) GetListrikMonitoringDataHarian() ([]entities.MonitoringData, error) {
	var monitoringDataHarian []entities.MonitoringData
	if err := r.db.Table("monitoring_data_harian").Where("monitoring_name LIKE ?", "monitoring_listrik_%").Find(&monitoringDataHarian).Error; err != nil {
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
