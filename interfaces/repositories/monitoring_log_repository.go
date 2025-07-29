package repositories

import "smartbuilding/entities"

type MonitoringLogRepository interface {
	SaveMonitoringLog(data entities.MonitoringLog) (entities.MonitoringLog, error)
	BulkSaveMonitoringLogs(data []entities.MonitoringLog) error
	GetAirMonitoringData(id int) ([]entities.MonitoringLog, error)
	GetListrikMonitoringData(id int) ([]entities.MonitoringLog, error)
	FindAll() ([]entities.MonitoringLog, error)
	FindBySettingID(id int) ([]entities.MonitoringLog, error)
	DeleteByIDs(ids []uint) error
	Truncate() error
}
