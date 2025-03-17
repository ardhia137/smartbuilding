package repositories

import "smartbuilding/entities"

type MonitoringDataRepository interface {
	SaveMonitoringData(data entities.MonitoringData) (entities.MonitoringData, error)
	GetAirMonitoringData(id int) ([]entities.MonitoringData, error)
	GetAirMonitoringDataHarian(id int) ([]entities.MonitoringData, error)
	GetListrikMonitoringData(id int) ([]entities.MonitoringData, error)
	GetListrikMonitoringDataHarian(id int) ([]entities.MonitoringData, error)
	FindAll() ([]entities.MonitoringData, error)
	FindBySettingID(id int) ([]entities.MonitoringData, error)
	SaveHarianData(data entities.MonitoringData) (*entities.MonitoringData, error)
	Truncate() error
}
