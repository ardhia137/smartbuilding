package services

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
)

type DataTorenServiceImpl struct {
	dataTorenRepo repositories.DataTorenRepository
}

func NewDataTorenService(dataTorenRepo repositories.DataTorenRepository) services.DataTorenService {
	return &DataTorenServiceImpl{dataTorenRepo: dataTorenRepo}
}

func (s *DataTorenServiceImpl) CreateDataToren(request entities.CreateDataTorenRequest) (*entities.DataTorenResponse, error) {
	dataToren := entities.DataToren{
		MonitoringName: request.MonitoringName,
		KapasitasToren: request.KapasitasToren,
		IDSetting:      request.IDSetting,
	}

	createdDataToren, err := s.dataTorenRepo.Create(&dataToren)
	if err != nil {
		return nil, err
	}

	response := entities.DataTorenResponse{
		ID:             createdDataToren.ID,
		MonitoringName: createdDataToren.MonitoringName,
		KapasitasToren: createdDataToren.KapasitasToren,
		IDSetting:      createdDataToren.IDSetting,
	}

	return &response, nil
}

func (s *DataTorenServiceImpl) GetAllDataToren() ([]entities.DataTorenResponse, error) {
	dataTorenList, err := s.dataTorenRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []entities.DataTorenResponse
	for _, dataToren := range dataTorenList {
		response = append(response, entities.DataTorenResponse{
			ID:             dataToren.ID,
			MonitoringName: dataToren.MonitoringName,
			KapasitasToren: dataToren.KapasitasToren,
			IDSetting:      dataToren.IDSetting,
		})
	}

	return response, nil
}

func (s *DataTorenServiceImpl) GetDataTorenBySettingID(id int) ([]entities.DataTorenResponse, error) {
	dataTorenList, err := s.dataTorenRepo.FindBySettingID(id)
	if err != nil {
		return nil, err
	}

	var response []entities.DataTorenResponse
	for _, dataToren := range dataTorenList {
		response = append(response, entities.DataTorenResponse{
			ID:             dataToren.ID,
			MonitoringName: dataToren.MonitoringName,
			KapasitasToren: dataToren.KapasitasToren,
			IDSetting:      dataToren.IDSetting,
		})
	}

	return response, nil
}

func (s *DataTorenServiceImpl) GetDataTorenByID(id int) (*entities.DataTorenResponse, error) {
	dataToren, err := s.dataTorenRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := entities.DataTorenResponse{
		ID:             dataToren.ID,
		MonitoringName: dataToren.MonitoringName,
		KapasitasToren: dataToren.KapasitasToren,
		IDSetting:      dataToren.IDSetting,
	}

	return &response, nil
}

func (s *DataTorenServiceImpl) UpdateDataToren(id int, request entities.CreateDataTorenRequest) (*entities.DataTorenResponse, error) {
	dataToren, err := s.dataTorenRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	dataToren.MonitoringName = request.MonitoringName
	dataToren.KapasitasToren = request.KapasitasToren
	dataToren.IDSetting = request.IDSetting

	updatedDataToren, err := s.dataTorenRepo.Update(dataToren)
	if err != nil {
		return nil, err
	}

	response := entities.DataTorenResponse{
		ID:             updatedDataToren.ID,
		MonitoringName: updatedDataToren.MonitoringName,
		KapasitasToren: updatedDataToren.KapasitasToren,
		IDSetting:      updatedDataToren.IDSetting,
	}

	return &response, nil
}

func (s *DataTorenServiceImpl) DeleteDataToren(id int) error {
	return s.dataTorenRepo.Delete(id)
}
