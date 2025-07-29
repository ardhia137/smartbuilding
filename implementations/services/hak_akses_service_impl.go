package services

import (
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
)

type HakAksesServiceImpl struct {
	hakAksesRepo repositories.HakAksesRepository
}

func NewHakAksesService(hakAksesRepo repositories.HakAksesRepository) services.HakAksesService {
	return &HakAksesServiceImpl{hakAksesRepo: hakAksesRepo}
}

func (s *HakAksesServiceImpl) CreateHakAkses(request entities.CreateHakAksesRequest) (*entities.HakAksesResponse, error) {
	hakAkses := entities.HakAkses{
		UserId:   request.UserID,
		GedungID: request.GedungID,
	}

	createdHakAkses, err := s.hakAksesRepo.Create(&hakAkses)
	if err != nil {
		return nil, err
	}

	response := entities.HakAksesResponse{
		ID:       createdHakAkses.ID,
		UserID:   createdHakAkses.UserId,
		GedungID: createdHakAkses.GedungID,
	}

	return &response, nil
}

func (s *HakAksesServiceImpl) GetAllHakAkses() ([]entities.AllHakAksesResponse, error) {
	hakAksesList, err := s.hakAksesRepo.FindAll()
	if err != nil {
		return nil, err
	}
	var response []entities.AllHakAksesResponse
	for _, hakAkses := range hakAksesList {
		response = append(response, entities.AllHakAksesResponse{
			ID:         hakAkses.ID,
			NamaGedung: hakAkses.NamaGedung,
			Username:   hakAkses.Username,
			Email:      hakAkses.Email,
			Role:       hakAkses.Role,
			GedungID:   hakAkses.GedungID,
		})
	}

	return response, nil
}

func (s *HakAksesServiceImpl) GetHakAksesByGedungIDUser(id int, userId int) ([]entities.HakAksesResponse, error) {
	hakAksesList, err := s.hakAksesRepo.FindByGedungIDUser(id, userId)
	if err != nil {
		return nil, err
	}

	var response []entities.HakAksesResponse
	for _, hakAkses := range hakAksesList {
		response = append(response, entities.HakAksesResponse{
			ID:       hakAkses.ID,
			UserID:   hakAkses.UserId,
			GedungID: hakAkses.GedungID,
		})
	}

	return response, nil
}

func (s *HakAksesServiceImpl) GetHakAksesByUser(userId int) ([]entities.AllHakAksesResponse, error) {
	hakAksesList, err := s.hakAksesRepo.FindByGedungUser(userId)
	if err != nil {
		return nil, err
	}

	var response []entities.AllHakAksesResponse
	for _, hakAkses := range hakAksesList {
		response = append(response, entities.AllHakAksesResponse{
			ID:         hakAkses.ID,
			NamaGedung: hakAkses.NamaGedung,
			Username:   hakAkses.Username,
			Email:      hakAkses.Email,
			Role:       hakAkses.Role,
			GedungID:   hakAkses.GedungID,
		})
	}

	return response, nil
}

func (s *HakAksesServiceImpl) GetHakAksesByID(id int) (*entities.HakAksesResponse, error) {
	hakAkses, err := s.hakAksesRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := entities.HakAksesResponse{
		ID:       hakAkses.ID,
		UserID:   hakAkses.UserId,
		GedungID: hakAkses.GedungID,
	}

	return &response, nil
}

func (s *HakAksesServiceImpl) UpdateHakAkses(id int, request entities.CreateHakAksesRequest) (*entities.HakAksesResponse, error) {
	hakAkses, err := s.hakAksesRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	hakAkses.UserId = request.UserID
	hakAkses.GedungID = request.GedungID

	updatedHakAkses, err := s.hakAksesRepo.Update(hakAkses)
	if err != nil {
		return nil, err
	}

	response := entities.HakAksesResponse{
		ID:       updatedHakAkses.ID,
		UserID:   updatedHakAkses.UserId,
		GedungID: updatedHakAkses.GedungID,
	}

	return &response, nil
}

func (s *HakAksesServiceImpl) DeleteHakAkses(id int) error {
	err := s.hakAksesRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete hak akses: %w", err)
	}
	return nil
}
