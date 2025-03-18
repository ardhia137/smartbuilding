package services

import (
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
)

type PengelolaGedungServiceImpl struct {
	pengelolaGedungRepo repositories.PengelolaGedungRepository
}

func NewPengelolaGedungService(pengelolaGedungRepo repositories.PengelolaGedungRepository) services.PengelolaGedungService {
	return &PengelolaGedungServiceImpl{pengelolaGedungRepo: pengelolaGedungRepo}
}

func (s *PengelolaGedungServiceImpl) CreatePengelolaGedung(request entities.CreatePengelolaGedungRequest) (*entities.PengelolaGedungResponse, error) {
	pengelolaGedung := entities.PengelolaGedung{
		UserId:    request.UserID,
		SettingID: request.SettingID,
	}

	createdPengelolaGedung, err := s.pengelolaGedungRepo.Create(&pengelolaGedung)
	if err != nil {
		return nil, err
	}

	response := entities.PengelolaGedungResponse{
		ID:        createdPengelolaGedung.ID,
		UserID:    createdPengelolaGedung.UserId,
		SettingID: createdPengelolaGedung.SettingID,
	}

	return &response, nil
}

func (s *PengelolaGedungServiceImpl) GetAllPengelolaGedung() ([]entities.AllPengelolaGedungResponse, error) {
	pengelolaGedungList, err := s.pengelolaGedungRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var response []entities.AllPengelolaGedungResponse
	for _, pengelolaGedung := range pengelolaGedungList {
		response = append(response, entities.AllPengelolaGedungResponse{
			ID:         pengelolaGedung.ID,
			NamaGedung: pengelolaGedung.NamaGedung,
			Username:   pengelolaGedung.Username,
			Email:      pengelolaGedung.Email,
			Role:       pengelolaGedung.Role,
		})
	}
	fmt.Println(response)
	return response, nil
}

func (s *PengelolaGedungServiceImpl) GetPengelolaGedungBySettingIDUser(id int, userID int) ([]entities.PengelolaGedungResponse, error) {

	pengelolaGedungList, err := s.pengelolaGedungRepo.FindBySettingIDUser(id, userID)
	if err != nil {
		return nil, err
	}

	var response []entities.PengelolaGedungResponse
	for _, pengelolaGedung := range pengelolaGedungList {
		response = append(response, entities.PengelolaGedungResponse{
			ID:        pengelolaGedung.ID,
			UserID:    pengelolaGedung.UserId,
			SettingID: pengelolaGedung.SettingID,
		})
	}

	return response, nil
}
func (s *PengelolaGedungServiceImpl) GetPengelolaGedungByUser(userID int) ([]entities.AllPengelolaGedungResponse, error) {

	pengelolaGedungList, err := s.pengelolaGedungRepo.FindBySettingUser(userID)
	if err != nil {
		return nil, err
	}

	var response []entities.AllPengelolaGedungResponse
	for _, pengelolaGedung := range pengelolaGedungList {
		response = append(response, entities.AllPengelolaGedungResponse{
			ID:         pengelolaGedung.ID,
			NamaGedung: pengelolaGedung.NamaGedung,
			Username:   pengelolaGedung.Username,
			Email:      pengelolaGedung.Email,
			Role:       pengelolaGedung.Role,
		})
	}

	return response, nil
}

func (s *PengelolaGedungServiceImpl) GetPengelolaGedungByID(id int) (*entities.PengelolaGedungResponse, error) {
	pengelolaGedung, err := s.pengelolaGedungRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := entities.PengelolaGedungResponse{
		ID:        pengelolaGedung.ID,
		UserID:    pengelolaGedung.UserId,
		SettingID: pengelolaGedung.SettingID,
	}

	return &response, nil
}

func (s *PengelolaGedungServiceImpl) UpdatePengelolaGedung(id int, request entities.CreatePengelolaGedungRequest) (*entities.PengelolaGedungResponse, error) {
	pengelolaGedung, err := s.pengelolaGedungRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	pengelolaGedung.UserId = request.UserID
	pengelolaGedung.SettingID = request.SettingID

	updatedPengelolaGedung, err := s.pengelolaGedungRepo.Update(pengelolaGedung)
	if err != nil {
		return nil, err
	}

	response := entities.PengelolaGedungResponse{
		ID:        updatedPengelolaGedung.ID,
		UserID:    updatedPengelolaGedung.UserId,
		SettingID: updatedPengelolaGedung.SettingID,
	}

	return &response, nil
}

func (s *PengelolaGedungServiceImpl) DeletePengelolaGedung(id int) error {
	return s.pengelolaGedungRepo.Delete(id)
}
