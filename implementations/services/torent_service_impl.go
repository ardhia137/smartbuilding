package services

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
)

type TorentServiceImpl struct {
	torentRepo repositories.TorentRepository
}

func NewTorentService(torentRepo repositories.TorentRepository) services.TorentService {
	return &TorentServiceImpl{torentRepo: torentRepo}
}

func (s *TorentServiceImpl) CreateTorent(request entities.CreateTorentRequest) (*entities.TorentResponse, error) {
	torent := entities.Torent{
		MonitoringName: request.MonitoringName,
		KapasitasToren: request.KapasitasToren,
		IDGedung:       request.IDGedung,
	}

	createdTorent, err := s.torentRepo.Create(&torent)
	if err != nil {
		return nil, err
	}

	response := entities.TorentResponse{
		ID:             createdTorent.ID,
		MonitoringName: createdTorent.MonitoringName,
		KapasitasToren: createdTorent.KapasitasToren,
		IDGedung:       createdTorent.IDGedung,
	}

	return &response, nil
}

func (s *TorentServiceImpl) GetAllTorent() ([]entities.TorentResponse, error) {
	torentList, err := s.torentRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []entities.TorentResponse
	for _, torent := range torentList {
		response = append(response, entities.TorentResponse{
			ID:             torent.ID,
			MonitoringName: torent.MonitoringName,
			KapasitasToren: torent.KapasitasToren,
			IDGedung:       torent.IDGedung,
		})
	}

	return response, nil
}

func (s *TorentServiceImpl) GetTorentByGedungID(id int) ([]entities.TorentResponse, error) {
	torentList, err := s.torentRepo.FindByGedungID(id)
	if err != nil {
		return nil, err
	}

	var response []entities.TorentResponse
	for _, torent := range torentList {
		response = append(response, entities.TorentResponse{
			ID:             torent.ID,
			MonitoringName: torent.MonitoringName,
			KapasitasToren: torent.KapasitasToren,
			IDGedung:       torent.IDGedung,
		})
	}

	return response, nil
}

func (s *TorentServiceImpl) GetTorentByID(id int) (*entities.TorentResponse, error) {
	torent, err := s.torentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := entities.TorentResponse{
		ID:             torent.ID,
		MonitoringName: torent.MonitoringName,
		KapasitasToren: torent.KapasitasToren,
		IDGedung:       torent.IDGedung,
	}

	return &response, nil
}

func (s *TorentServiceImpl) UpdateTorent(id int, request entities.CreateTorentRequest) (*entities.TorentResponse, error) {
	torent, err := s.torentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	torent.MonitoringName = request.MonitoringName
	torent.KapasitasToren = request.KapasitasToren
	torent.IDGedung = request.IDGedung

	updatedTorent, err := s.torentRepo.Update(torent)
	if err != nil {
		return nil, err
	}

	response := entities.TorentResponse{
		ID:             updatedTorent.ID,
		MonitoringName: updatedTorent.MonitoringName,
		KapasitasToren: updatedTorent.KapasitasToren,
		IDGedung:       updatedTorent.IDGedung,
	}

	return &response, nil
}

func (s *TorentServiceImpl) DeleteTorent(id int) error {
	return s.torentRepo.Delete(id)
}
