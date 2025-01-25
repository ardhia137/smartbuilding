package repositories

import "smartbuilding/entities"

type MonitoringDataRepository interface {
	SaveMonitoringData(data entities.MonitoringData) (entities.MonitoringData, error)
	GetAirMonitoringData() ([]entities.MonitoringData, error)
	GetAirMonitoringDataHarian() ([]entities.MonitoringData, error)
	GetListrikMonitoringData() ([]entities.MonitoringData, error)
	GetListrikMonitoringDataHarian() ([]entities.MonitoringData, error)
	FindAll() ([]entities.MonitoringData, error)
	SaveHarianData(data entities.MonitoringData) (*entities.MonitoringData, error)
	Truncate() error
}
