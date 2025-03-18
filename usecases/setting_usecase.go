package usecases

import (
	"smartbuilding/entities"
)

type SettingUseCase interface {
	CreateSetting(request entities.CreateSettingRequest) (*entities.SettingResponseCreate, error)
	GetAllSetting(role string, userID uint) ([]entities.SettingResponse, error) // Perubahan di sini
	GetAllCornJobs() ([]entities.SettingResponse, error)
	GetSettingByID(id int) (*entities.SettingResponse, error)
	UpdateSetting(id int, request entities.CreateSettingRequest) (*entities.SettingResponse, error)
	DeleteSetting(id int) error
}
