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
		// Gagal menghitung jumlah penyewa
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	fmt.Println(jumlahPenyewa)
	// Jika kapasitas kamar sudah penuh, ubah status kamar menjadi tidak tersedia
	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia" // Misalnya "Tidak Tersedia"
		// Update status kamar
		_, err := s.kamarRepository.Update(kamar.ID, kamar)
		if err != nil {
			return entities.PenyewaKamarResponse{}, utils.ErrInternal
		}

		// Kamar sudah penuh, langsung return response
		return entities.PenyewaKamarResponse{}, utils.ErrKamar
	}

	// Parse TanggalMulai
	tanggalMulai, err := time.Parse("2006-01-02", request.TanggalMulai)
	if err != nil {
		// Gagal parse tanggal
		return entities.PenyewaKamarResponse{}, utils.ErrBadRequest
	}

	// Jika NPM belum terdaftar, lanjutkan untuk membuat pengguna baru
	penyewaKamar = entities.PenyewaKamar{
		ID:           request.ID,
		NPM:          request.NPM,
		KamarID:      request.KamarID,
		TanggalMulai: tanggalMulai,
		Status:       request.Status,
	}

	// Simpan data penyewa kamar ke repository
	createdPenyewaKamar, err := s.penyewaKamarRepository.Create(penyewaKamar)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	jumlahPenyewa = +1
	fmt.Println(jumlahPenyewa)
	// Pastikan status kamar diperbarui jika jumlah penyewa berkurang
	// Update status kamar jika jumlah penyewa berkurang dan kapasitas belum tercapai
	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia" // Kamar bisa menjadi tersedia jika ada tempat
		_, err := s.kamarRepository.Update(kamar.ID, kamar)
		fmt.Println(jumlahPenyewa)
		if err != nil {
			return entities.PenyewaKamarResponse{}, utils.ErrInternal
		}
	}

	// Mendapatkan data terkait Mahasiswa dan User
	mahasiswa, err := s.mahasiswaRepository.FindByID(request.NPM)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	user, err := s.userRepository.FindByID(mahasiswa.UserID)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	// Gunakan mapPenyewaKamarToResponse untuk membentuk response yang lebih terstruktur
	return mapPenyewaKamarToResponse(createdPenyewaKamar, user, mahasiswa, kamar), nil
}
func (s *penyewaKamarServiceImpl) UpdatePenyewaKamar(id uint, request entities.UpdatePenyewaKamarRequest) (entities.PenyewaKamarResponse, error) {
	// Mencari penyewa kamar berdasarkan ID
	penyewaKamar, err := s.penyewaKamarRepository.FindByID(id)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrNotFound
	}

	// Memeriksa ketersediaan kamar
	kamar, err := s.kamarRepository.FindByID(request.KamarID)
	if err != nil || kamar.ID == 0 {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	// Parse TanggalMulai dan TanggalKeluar
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

	// Memperbarui data penyewa kamar
	penyewaKamar.NPM = request.NPM
	penyewaKamar.KamarID = request.KamarID
	penyewaKamar.TanggalMulai = tanggalMulai
	penyewaKamar.TanggalKeluar = tanggalKeluar
	penyewaKamar.Status = request.Status

	// Menyimpan perubahan ke database
	updatedPenyewaKamar, err := s.penyewaKamarRepository.Update(id, penyewaKamar)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	// Menghitung jumlah penyewa yang ada di kamar
	jumlahPenyewa, err := s.penyewaKamarRepository.CountAktifByKamarID(request.KamarID)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}
	fmt.Println(jumlahPenyewa)
	fmt.Println(kamar.Kapasitas)
	// Tentukan status kamar berdasarkan jumlah penyewa dan kapasitas
	if jumlahPenyewa >= int64(kamar.Kapasitas) {
		kamar.Status = "tidak tersedia"
	} else {
		kamar.Status = "tersedia"
	}

	// Update status kamar setelah perubahan penyewa
	_, err = s.kamarRepository.Update(kamar.ID, kamar)
	if err != nil {
		return entities.PenyewaKamarResponse{}, utils.ErrInternal
	}

	// Mendapatkan data User, Mahasiswa, dan Kamar terkait untuk response
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

	// Kembalikan response yang terstruktur
	return mapPenyewaKamarToResponse(updatedPenyewaKamar, user, mahasiswa, kamar), nil
}

// DeletePenyewaKamar menghapus pengguna berdasarkan ID
func (s *penyewaKamarServiceImpl) DeletePenyewaKamar(id uint) error {
	// Mencari pengguna berdasarkan ID
	penyewakamar, err := s.penyewaKamarRepository.FindByID(id)
	if err != nil {
		return utils.ErrNotFound // Jika pengguna tidak ditemukan
	}
	kamar, err := s.kamarRepository.FindByID(penyewakamar.KamarID)
	if err != nil || kamar.ID == 0 {
		return utils.ErrInternal
	}
	// Menghapus pengguna dari database
	err = s.penyewaKamarRepository.Delete(id)
	if err != nil {
		return utils.ErrInternal // Jika terjadi error saat penghapusan
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

	// Update status kamar setelah perubahan penyewa
	_, err = s.kamarRepository.Update(kamar.ID, kamar)
	if err != nil {
		return utils.ErrInternal
	}

	return nil // Penghapusan berhasil
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
		TanggalKeluar: tanggalKeluarStr, // Akan kosong jika TanggalKeluar nil
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
