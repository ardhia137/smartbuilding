package services

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
)

// kamarServiceImpl adalah struktur implementasi KamarService
type kamarServiceImpl struct {
	kamarRepository repositories.KamarRepository
}

// NewKamarService membuat instance dari KamarService
func NewKamarService(kamarRepo repositories.KamarRepository) services.KamarService {
	return &kamarServiceImpl{kamarRepo}
}

// GetAllKamars mengembalikan semua pengguna
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

// GetKamarByID mengembalikan informasi pengguna berdasarkan ID
func (s *kamarServiceImpl) GetKamarByID(id uint) (entities.KamarResponse, error) {
	kamar, err := s.kamarRepository.FindByID(id)
	if err != nil {
		return entities.KamarResponse{}, utils.ErrNotFound
	}

	return entities.KamarResponse{ID: kamar.ID, NoKamar: kamar.NoKamar, Lantai: kamar.Lantai, Kapasitas: kamar.Kapasitas, Status: kamar.Status}, nil
}

// CreateKamar membuat pengguna baru
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

// UpdateKamar memperbarui informasi pengguna berdasarkan ID
func (s *kamarServiceImpl) UpdateKamar(id uint, request entities.CreateKamarRequest) (entities.KamarResponse, error) {
	// Mencari pengguna berdasarkan ID
	kamar, err := s.kamarRepository.FindByID(id)
	if err != nil {
		return entities.KamarResponse{}, utils.ErrNotFound
	}

	// Memperbarui data pengguna
	kamar.NoKamar = request.NoKamar
	kamar.Lantai = request.Lantai
	kamar.Kapasitas = request.Kapasitas
	kamar.Status = request.Status

	// Menyimpan perubahan ke database dengan id dan kamar
	updatedKamar, err := s.kamarRepository.Update(id, kamar) // Memasukkan id dan kamar
	if err != nil {
		return entities.KamarResponse{}, utils.ErrInternal
	}

	// Mengembalikan kamar yang telah diperbarui
	return entities.KamarResponse{
		ID:        updatedKamar.ID,
		NoKamar:   updatedKamar.NoKamar,
		Lantai:    updatedKamar.Lantai,
		Kapasitas: updatedKamar.Kapasitas,
		Status:    updatedKamar.Status,
	}, nil
}

// DeleteKamar menghapus pengguna berdasarkan ID
func (s *kamarServiceImpl) DeleteKamar(id uint) error {
	// Mencari pengguna berdasarkan ID
	_, err := s.kamarRepository.FindByID(id)
	if err != nil {
		return utils.ErrNotFound // Jika pengguna tidak ditemukan
	}

	// Menghapus pengguna dari database
	err = s.kamarRepository.Delete(id)
	if err != nil {
		return utils.ErrInternal // Jika terjadi error saat penghapusan
	}

	return nil // Penghapusan berhasil
}
