package services

import (
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
	"time"
)

type penyewaKamarServiceImpl struct {
	penyewaKamarRepository repositories.PenyewaKamarRepository
	kamarRepository        repositories.KamarRepository
	userRepository         repositories.UserRepository
	mahasiswaRepository    repositories.MahasiswaRepository
}

func NewPenyewaKamarService(
	penyewaKamarRepo repositories.PenyewaKamarRepository,
	kamarRepo repositories.KamarRepository,
	userRepo repositories.UserRepository,
	mahasiswaRepo repositories.MahasiswaRepository,
) services.PenyewaKamarService {
	return &penyewaKamarServiceImpl{
		penyewaKamarRepo,
		kamarRepo,
		userRepo,
		mahasiswaRepo,
	}
}

func (s *penyewaKamarServiceImpl) GetAllPenyewaKamar() ([]entities.PenyewaKamarResponse, error) {
	penyewaKamars, err := s.penyewaKamarRepository.FindAll()
	if err != nil {
		return nil, utils.ErrInternal
	}

	var penyewaKamarResponses []entities.PenyewaKamarResponse
	for _, penyewaKamar := range penyewaKamars {
		mahasiswa, _ := s.mahasiswaRepository.FindByID(penyewaKamar.NPM)
		user, _ := s.userRepository.FindByID(mahasiswa.UserID)
		kamar, _ := s.kamarRepository.FindByID(penyewaKamar.KamarID)
		response := mapPenyewaKamarToResponse(penyewaKamar, user, mahasiswa, kamar)
		penyewaKamarResponses = append(penyewaKamarResponses, response)
	}
	return penyewaKamarResponses, nil
}

func (s *penyewaKamarServiceImpl) GetPenyewaKamarByID(id uint) (entities.PenyewaKamarResponse, error) {
	penyewaKamar, err := s.penyewaKamarRepository.FindByID(id)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrNotFound
	}

	mahasiswa, _ := s.mahasiswaRepository.FindByID(penyewaKamar.NPM)
	user, _ := s.userRepository.FindByID(mahasiswa.UserID)
	kamar, _ := s.kamarRepository.FindByID(penyewaKamar.KamarID)

	return mapPenyewaKamarToResponse(penyewaKamar, user, mahasiswa, kamar), nil
}

func (s *penyewaKamarServiceImpl) FindByNPM(npm uint) (entities.PenyewaKamarResponse, error) {
	penyewaKamar, err := s.penyewaKamarRepository.FindByNPM(npm)
	if err != nil {
		return entities.PenyewaKamarResponse{}, err
	}
	mahasiswa, _ := s.mahasiswaRepository.FindByID(penyewaKamar.NPM)
	user, _ := s.userRepository.FindByID(mahasiswa.UserID)
	kamar, _ := s.kamarRepository.FindByID(penyewaKamar.KamarID)
	var tanggalKeluar time.Time
	if penyewaKamar.TanggalKeluar != nil {
		tanggalKeluar = *penyewaKamar.TanggalKeluar
	}

	if tanggalKeluar.IsZero() {
		return mapPenyewaKamarToResponse(penyewaKamar, user, mahasiswa, kamar), nil
	}
	return mapPenyewaKamarToResponse(penyewaKamar, user, mahasiswa, kamar), nil
}

func (s *penyewaKamarServiceImpl) CekKetersediaanKamar(kamarID uint) (bool, error) {
	kamar, err := s.kamarRepository.FindByID(kamarID)
	if err != nil || kamar.ID == 0 {
		return false, utils.ErrNotFound
	}

	jumlahPenyewaAktif, err := s.penyewaKamarRepository.CountAktifByKamarID(kamarID)
	if err != nil {
		return false, utils.ErrInternal
	}

	if jumlahPenyewaAktif >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia"
		_, err := s.kamarRepository.Update(kamar.ID, kamar)
		if err != nil {
			return false, utils.ErrInternal
		}
		return false, nil
	}

	kamar.Status = "tersedia"
	_, err = s.kamarRepository.Update(kamar.ID, kamar)
	if err != nil {
		return false, utils.ErrInternal
	}

	return true, nil
}

func (s *penyewaKamarServiceImpl) CreatePenyewaKamar(request entities.CreatePenyewaKamarRequest) (entities.PenyewaKamarResponse, error) {
	penyewaKamar, err := s.penyewaKamarRepository.FindByNPM(request.NPM)
	if err == nil && penyewaKamar.NPM != 0 {
		return entities.PenyewaKamarResponse{}, utils.ErrNPMAlreadyExists
	}
	kamar, err := s.kamarRepository.FindByID(request.KamarID)
	if err != nil || kamar.ID == 0 {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	jumlahPenyewa, err := s.penyewaKamarRepository.CountAktifByKamarID(request.KamarID)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	fmt.Println(jumlahPenyewa)
	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia"
		_, err := s.kamarRepository.Update(kamar.ID, kamar)
		if err != nil {
			return entities.PenyewaKamarResponse{}, utils.ErrInternal
		}
		return entities.PenyewaKamarResponse{}, utils.ErrKamar
	}

	tanggalMulai, err := time.Parse("2006-01-02", request.TanggalMulai)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrBadRequest
	}

	penyewaKamar = entities.PenyewaKamar{
		ID:           request.ID,
		NPM:          request.NPM,
		KamarID:      request.KamarID,
		TanggalMulai: tanggalMulai,
		Status:       request.Status,
	}

	createdPenyewaKamar, err := s.penyewaKamarRepository.Create(penyewaKamar)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	jumlahPenyewa = +1
	fmt.Println(jumlahPenyewa)
	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia"
		_, err := s.kamarRepository.Update(kamar.ID, kamar)
		fmt.Println(jumlahPenyewa)
		if err != nil {
			return entities.PenyewaKamarResponse{}, utils.ErrInternal
		}
	}

	mahasiswa, err := s.mahasiswaRepository.FindByID(request.NPM)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	user, err := s.userRepository.FindByID(mahasiswa.UserID)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	return mapPenyewaKamarToResponse(createdPenyewaKamar, user, mahasiswa, kamar), nil
}

func (s *penyewaKamarServiceImpl) UpdatePenyewaKamar(id uint, request entities.UpdatePenyewaKamarRequest) (entities.PenyewaKamarResponse, error) {
	penyewaKamar, err := s.penyewaKamarRepository.FindByID(id)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrNotFound
	}

	kamar, err := s.kamarRepository.FindByID(request.KamarID)
	if err != nil || kamar.ID == 0 {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	tanggalMulai, err := time.Parse("2006-01-02", request.TanggalMulai)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrBadRequest
	}

	var tanggalKeluar *time.Time
	if request.TanggalKeluar != "" {
		parsedTanggalKeluar, err := time.Parse("2006-01-02", request.TanggalKeluar)
		if err != nil {
			return entities.PenyewaKamarResponse{}, utils.ErrBadRequest
		}
		tanggalKeluar = &parsedTanggalKeluar
	}

	penyewaKamar.NPM = request.NPM
	penyewaKamar.KamarID = request.KamarID
	penyewaKamar.TanggalMulai = tanggalMulai
	penyewaKamar.TanggalKeluar = tanggalKeluar
	penyewaKamar.Status = request.Status

	updatedPenyewaKamar, err := s.penyewaKamarRepository.Update(id, penyewaKamar)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	jumlahPenyewa, err := s.penyewaKamarRepository.CountAktifByKamarID(request.KamarID)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	fmt.Println(jumlahPenyewa)
	fmt.Println(kamar.Kapasitas)
	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia"
	} else {
		kamar.Status = "tersedia"
	}

	_, err = s.kamarRepository.Update(kamar.ID, kamar)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	mahasiswa, err := s.mahasiswaRepository.FindByID(updatedPenyewaKamar.NPM)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	user, err := s.userRepository.FindByID(mahasiswa.UserID)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	kamar, err = s.kamarRepository.FindByID(updatedPenyewaKamar.KamarID)
	if err != nil || kamar.ID == 0 {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	return mapPenyewaKamarToResponse(updatedPenyewaKamar, user, mahasiswa, kamar), nil
}

func (s *penyewaKamarServiceImpl) DeletePenyewaKamar(id uint) error {
	penyewakamar, err := s.penyewaKamarRepository.FindByID(id)
	if err != nil {
		return utils.ErrNotFound
	}
	kamar, err := s.kamarRepository.FindByID(penyewakamar.KamarID)
	if err != nil || kamar.ID == 0 {
		return utils.ErrInternal
	}
	err = s.penyewaKamarRepository.Delete(id)
	if err != nil {
		return utils.ErrInternal
	}
	jumlahPenyewa, err := s.penyewaKamarRepository.CountAktifByKamarID(penyewakamar.KamarID)
	if err != nil {
		return utils.ErrInternal
	}

	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia"
	} else {
		kamar.Status = "tersedia"
	}

	_, err = s.kamarRepository.Update(kamar.ID, kamar)
	if err != nil {
		return utils.ErrInternal
	}

	return nil
}

func mapPenyewaKamarToResponse(penyewaKamar entities.PenyewaKamar, user entities.User, mahasiswa entities.Mahasiswa, kamar entities.Kamar) entities.PenyewaKamarResponse {
	var tanggalKeluarStr string
	if penyewaKamar.TanggalKeluar != nil {
		tanggalKeluarStr = penyewaKamar.TanggalKeluar.Format("2006-01-02")
	}

	return entities.PenyewaKamarResponse{
		ID:            penyewaKamar.ID,
		NPM:           penyewaKamar.NPM,
		KamarID:       penyewaKamar.KamarID,
		TanggalMulai:  penyewaKamar.TanggalMulai.Format("2006-01-02"),
		TanggalKeluar: tanggalKeluarStr,
		Status:        penyewaKamar.Status,
		Mahasiswa:     mapMahasiswaToResponse(mahasiswa, user),
		Kamar: entities.KamarResponse{
			ID:        kamar.ID,
			NoKamar:   kamar.NoKamar,
			Lantai:    kamar.Lantai,
			Kapasitas: kamar.Kapasitas,
			Status:    kamar.Status,
		},
	}
}
