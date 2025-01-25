package services

import "smartbuilding/entities"

type MonitoringDataService interface {
	SaveMonitoringData(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error)
	GetAirMonitoringData() ([]entities.GetAirDataResponse, error)
	GetListrikMonitoringData() (entities.GetListrikDataResponse, error)
}
