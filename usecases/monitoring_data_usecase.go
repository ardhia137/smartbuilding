package usecases

import "smartbuilding/entities"

type MonitoringDataUseCase interface {
	SaveMonitoringData(request entities.CreateMonitoringDataRequest) (entities.MonitoringDataResponse, error)
	GetAirMonitoringData() ([]entities.GetAirDataResponse, error)
	GetListrikMonitoringData() (entities.GetListrikDataResponse, error)
}
