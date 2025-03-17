package usecases

import "smartbuilding/entities"

type DataTorenUseCase interface {
	CreateDataToren(request entities.CreateDataTorenRequest) (*entities.DataTorenResponse, error)
	GetAllDataToren() ([]entities.DataTorenResponse, error)
	GetDataTorenByID(id int) (*entities.DataTorenResponse, error)
	GetDataTorenBySettingID(id int) ([]entities.DataTorenResponse, error)
	UpdateDataToren(id int, request entities.CreateDataTorenRequest) (*entities.DataTorenResponse, error)
	DeleteDataToren(id int) error
}
