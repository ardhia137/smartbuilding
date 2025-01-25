package usecases

import "smartbuilding/entities"

type SettingUseCase interface {
	CreateSetting(request entities.CreateSettingRequest) (*entities.SettingResponse, error)
	GetAllSetting() ([]entities.SettingResponse, error)
	GetSettingByID(id int) (*entities.SettingResponse, error)
	UpdateSetting(id int, request entities.CreateSettingRequest) (*entities.SettingResponse, error)
	DeleteSetting(id int) error
}
