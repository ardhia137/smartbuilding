package repositories

import (
	"smartbuilding/entities"

	"gorm.io/gorm"
)

type MonitoringLogRepositoryImpl struct {
	db *gorm.DB
}

func NewMonitoringLogRepository(db *gorm.DB) *MonitoringLogRepositoryImpl {
	return &MonitoringLogRepositoryImpl{db: db}
}

func (r *MonitoringLogRepositoryImpl) SaveMonitoringLog(monitoringLog entities.MonitoringLog) (entities.MonitoringLog, error) {
	if err := r.db.Create(&monitoringLog).Error; err != nil {
		return entities.MonitoringLog{}, err
	}
	return monitoringLog, nil
}

func (r *MonitoringLogRepositoryImpl) BulkSaveMonitoringLogs(data []entities.MonitoringLog) error {
	if len(data) == 0 {
		return nil
	}

	// GORM bulk insert - much faster than individual inserts
	if err := r.db.CreateInBatches(&data, 100).Error; err != nil {
		return err
	}
	return nil
}

func (r *MonitoringLogRepositoryImpl) GetAirMonitoringData(id int) ([]entities.MonitoringLog, error) {
	var monitoringLog []entities.MonitoringLog
	if err := r.db.Where("monitoring_name LIKE ?", "monitoring_air%").Where("id_gedung = ?", id).Find(&monitoringLog).Error; err != nil {
		return nil, err
	}
	return monitoringLog, nil
}

func (r *MonitoringLogRepositoryImpl) GetListrikMonitoringData(id int) ([]entities.MonitoringLog, error) {
	var monitoringLog []entities.MonitoringLog
	if err := r.db.Where("monitoring_name LIKE ?", "monitoring_listrik_%").Where("id_gedung = ?", id).Find(&monitoringLog).Error; err != nil {
		return nil, err
	}
	return monitoringLog, nil
}

func (r *MonitoringLogRepositoryImpl) FindAll() ([]entities.MonitoringLog, error) {
	var data []entities.MonitoringLog
	if err := r.db.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MonitoringLogRepositoryImpl) FindBySettingID(id int) ([]entities.MonitoringLog, error) {
	var data []entities.MonitoringLog
	if err := r.db.Find(&data).Where("id_gedung = ?", id).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (r *MonitoringLogRepositoryImpl) DeleteByIDs(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := r.db.Where("id IN ?", ids).Delete(&entities.MonitoringLog{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *MonitoringLogRepositoryImpl) Truncate() error {
	if err := r.db.Exec("TRUNCATE TABLE monitoring_logs").Error; err != nil {
		return err
	}
	return nil
}
