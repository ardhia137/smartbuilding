package usecases

import "smartbuilding/entities"

type MonitoringLogUseCase interface {
	SaveMonitoringLog(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error)
	GetAirMonitoringData(id int) ([]entities.GetAirDataResponse, error)
	GetListrikMonitoringData(id int) (entities.GetListrikDataResponse, error)
}
