package services

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
)

type SettingServiceImpl struct {
	haosRepo repositories.SettingRepository
}

func NewSettingService(haosRepo repositories.SettingRepository) services.SettingService {
	return &SettingServiceImpl{haosRepo: haosRepo}
}

func (s *SettingServiceImpl) CreateSetting(request entities.CreateSettingRequest) (*entities.SettingResponse, error) {
	haos := entities.Setting{
		HaosURL:   request.HaosURL,
		HaosToken: request.HaosToken,
		Scheduler: request.Scheduler,
	}

	createdSetting, err := s.haosRepo.Create(&haos)
	if err != nil {
		return nil, err
	}

	response := entities.SettingResponse{
		ID:        createdSetting.ID,
		HaosURL:   createdSetting.HaosURL,
		HaosToken: createdSetting.HaosToken,
		Scheduler: createdSetting.Scheduler,
	}

	return &response, nil
}

func (s *SettingServiceImpl) GetAllSetting() ([]entities.SettingResponse, error) {
	haosList, err := s.haosRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []entities.SettingResponse
	for _, haos := range haosList {
		response = append(response, entities.SettingResponse{
			ID:        haos.ID,
			HaosURL:   haos.HaosURL,
			HaosToken: haos.HaosToken,
			Scheduler: haos.Scheduler,
		})
	}

	return response, nil
}

func (s *SettingServiceImpl) GetSettingByID(id int) (*entities.SettingResponse, error) {
	haos, err := s.haosRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := entities.SettingResponse{
		ID:        haos.ID,
		HaosURL:   haos.HaosURL,
		HaosToken: haos.HaosToken,
		Scheduler: haos.Scheduler,
	}

	return &response, nil
}

func (s *SettingServiceImpl) UpdateSetting(id int, request entities.CreateSettingRequest) (*entities.SettingResponse, error) {
	haos, err := s.haosRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	haos.HaosURL = request.HaosURL
	haos.HaosToken = request.HaosToken
	haos.Scheduler = request.Scheduler

	updatedSetting, err := s.haosRepo.Update(haos)
	if err != nil {
		return nil, err
	}

	response := entities.SettingResponse{
		ID:        updatedSetting.ID,
		HaosURL:   updatedSetting.HaosURL,
		HaosToken: updatedSetting.HaosToken,
		Scheduler: updatedSetting.Scheduler,
	}

	return &response, nil
}

func (s *SettingServiceImpl) DeleteSetting(id int) error {
	return s.haosRepo.Delete(id)
}
