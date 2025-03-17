package services

import "smartbuilding/entities"

type SettingService interface {
	CreateSetting(request entities.CreateSettingRequest) (*entities.SettingResponseCreate, error)
	GetAllSetting() ([]entities.SettingResponse, error)
	GetSettingByID(id int) (*entities.SettingResponse, error)
	UpdateSetting(id int, request entities.CreateSettingRequest) (*entities.SettingResponse, error)
	DeleteSetting(id int) error
}
