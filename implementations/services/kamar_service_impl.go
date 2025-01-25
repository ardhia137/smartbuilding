package services

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
)

type kamarServiceImpl struct {
	kamarRepository repositories.KamarRepository
}

func NewKamarService(kamarRepo repositories.KamarRepository) services.KamarService {
	return &kamarServiceImpl{kamarRepo}
}

func (s *kamarServiceImpl) GetAllKamar() ([]entities.KamarResponse, error) {
	kamars, err := s.kamarRepository.FindAll()
	if err != nil {
		return nil, utils.ErrInternal
	}

	var kamarResponses []entities.KamarResponse
	for _, kamar := range kamars {
		kamarResponses = append(kamarResponses, entities.KamarResponse{
			ID:        kamar.ID,
			NoKamar:   kamar.NoKamar,
			Lantai:    kamar.Lantai,
			Kapasitas: kamar.Kapasitas,
			Status:    kamar.Status,
		})
	}
	return kamarResponses, nil
}

func (s *kamarServiceImpl) GetKamarByID(id uint) (entities.KamarResponse, error) {
	kamar, err := s.kamarRepository.FindByID(id)
	if err != nil {
		return entities.KamarResponse{}, utils.ErrNotFound
	}

	return entities.KamarResponse{ID: kamar.ID, NoKamar: kamar.NoKamar, Lantai: kamar.Lantai, Kapasitas: kamar.Kapasitas, Status: kamar.Status}, nil
}

func (s *kamarServiceImpl) CreateKamar(request entities.CreateKamarRequest) (entities.KamarResponse, error) {
	kamar := entities.Kamar{
		ID:        request.ID,
		NoKamar:   request.NoKamar,
		Lantai:    request.Lantai,
		Kapasitas: request.Kapasitas,
		Status:    request.Status,
	}
	createdKamar, err := s.kamarRepository.Create(kamar)
	if err != nil {
		return entities.KamarResponse{}, utils.ErrInternal
	}

	return entities.KamarResponse{ID: createdKamar.ID, NoKamar: createdKamar.NoKamar, Lantai: createdKamar.Lantai, Kapasitas: kamar.Kapasitas, Status: createdKamar.Status}, nil
}

func (s *kamarServiceImpl) UpdateKamar(id uint, request entities.CreateKamarRequest) (entities.KamarResponse, error) {
	kamar, err := s.kamarRepository.FindByID(id)
	if err != nil {
		return entities.KamarResponse{}, utils.ErrNotFound
	}

	kamar.NoKamar = request.NoKamar
	kamar.Lantai = request.Lantai
	kamar.Kapasitas = request.Kapasitas
	kamar.Status = request.Status

	updatedKamar, err := s.kamarRepository.Update(id, kamar)
	if err != nil {
		return entities.KamarResponse{}, utils.ErrInternal
	}

	return entities.KamarResponse{
		ID:        updatedKamar.ID,
		NoKamar:   updatedKamar.NoKamar,
		Lantai:    updatedKamar.Lantai,
		Kapasitas: updatedKamar.Kapasitas,
		Status:    updatedKamar.Status,
	}, nil
}

func (s *kamarServiceImpl) DeleteKamar(id uint) error {
	_, err := s.kamarRepository.FindByID(id)
	if err != nil {
		return utils.ErrNotFound
	}

	err = s.kamarRepository.Delete(id)
	if err != nil {
		return utils.ErrInternal
	}

	return nil
}
