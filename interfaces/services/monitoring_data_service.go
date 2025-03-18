package services

import "smartbuilding/entities"

type MonitoringDataService interface {
	SaveMonitoringData(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error)
	GetAirMonitoringData(id int) ([]entities.GetAirDataResponse, error)
	GetListrikMonitoringData(id int) (entities.GetListrikDataResponse, error)
	//GetAirMonitoringDataCornJobs(id int) ([]entities.GetAirDataResponse, error)
	//GetListrikMonitoringDataCornJobs(id int) (entities.GetListrikDataResponse, error)
}
