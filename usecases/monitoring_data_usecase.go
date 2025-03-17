package usecases

import "smartbuilding/entities"

type MonitoringDataUseCase interface {
	SaveMonitoringData(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error)
	GetAirMonitoringData(id int) ([]entities.GetAirDataResponse, error)
	GetListrikMonitoringData(id int) (entities.GetListrikDataResponse, error)
}
