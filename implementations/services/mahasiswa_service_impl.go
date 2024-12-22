package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
	"time"
)

// mahasiswaServiceImpl adalah struktur implementasi MahasiswaService
type mahasiswaServiceImpl struct {
	mahasiswaRepository repositories.MahasiswaRepository
	userRepository      repositories.UserRepository
}

// NewMahasiswaService membuat instance dari MahasiswaService
func NewMahasiswaService(
	mahasiswaRepo repositories.MahasiswaRepository,
	userRepo repositories.UserRepository,
) services.MahasiswaService {
	return &mahasiswaServiceImpl{
		mahasiswaRepository: mahasiswaRepo,
		userRepository:      userRepo,
	}
}

// GetAllMahasiswa mengembalikan semua data mahasiswa
func (s *mahasiswaServiceImpl) GetAllMahasiswa() ([]entities.MahasiswaResponse, error) {
	mahasiswas, err := s.mahasiswaRepository.FindAll()
	if err != nil {
		return nil, utils.ErrInternal
	}

	var mahasiswaResponses []entities.MahasiswaResponse
	for _, mahasiswa := range mahasiswas {
		user, err := s.userRepository.FindByID(mahasiswa.UserID)
		if err != nil {
			return nil, utils.ErrInternal
		}
		mahasiswaResponses = append(mahasiswaResponses, mapMahasiswaToResponse(mahasiswa, user))
	}

	return mahasiswaResponses, nil
}

// GetMahasiswaByID mengembalikan informasi mahasiswa berdasarkan ID
func (s *mahasiswaServiceImpl) GetMahasiswaByID(NPM uint) (entities.MahasiswaResponse, error) {
	fmt.Print(NPM)
	mahasiswa, err := s.mahasiswaRepository.FindByID(NPM)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrNotFound
	}

	user, err := s.userRepository.FindByID(mahasiswa.UserID)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	return mapMahasiswaToResponse(mahasiswa, user), nil
}

// CreateMahasiswa membuat data mahasiswa baru
func (s *mahasiswaServiceImpl) CreateMahasiswa(request entities.CreateMahasiswaRequest) (entities.MahasiswaResponse, error) {
	// Buat data user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.User.Password), bcrypt.DefaultCost)
	fmt.Print(request.TanggalMasuk)
	user := entities.User{
		Username: request.User.Username,
		Email:    request.User.Email,
		Password: string(hashedPassword),
		Role:     "mahasiswa",
	}
	createdUser, err := s.userRepository.Create(user)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}
	tanggalMasuk, err := time.Parse("2006-01-02", request.TanggalMasuk)
	tanggalLahir, err := time.Parse("2006-01-02", request.TanggalMasuk)
	fmt.Print(request)
	// Buat data mahasiswa
	mahasiswa := entities.Mahasiswa{
		NPM:             request.NPM,
		Nama:            request.Nama,
		TanggalLahir:    tanggalLahir,
		Fakultas:        request.Fakultas,
		Jurusan:         request.Jurusan,
		TanggalMasuk:    tanggalMasuk,
		JenisKelamin:    request.JenisKelamin,
		StatusMahasiswa: request.StatusMahasiswa,
		UserID:          createdUser.ID,
	}
	createdMahasiswa, err := s.mahasiswaRepository.Create(mahasiswa)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	return mapMahasiswaToResponse(createdMahasiswa, createdUser), nil
}

// UpdateMahasiswa memperbarui data mahasiswa dan data pengguna terkait berdasarkan ID
func (s *mahasiswaServiceImpl) UpdateMahasiswa(NPM uint, request entities.UpdateMahasiswaRequest) (entities.MahasiswaResponse, error) {
	// Cari data mahasiswa
	mahasiswa, err := s.mahasiswaRepository.FindByID(NPM)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrNotFound
	}
	tanggalMasuk, err := time.Parse("2006-01-02", request.TanggalMasuk)
	tanggalLahir, err := time.Parse("2006-01-02", request.TanggalMasuk)
	fmt.Print(request.TanggalMasuk)
	fmt.Print(tanggalMasuk)
	// Perbarui data mahasiswa
	mahasiswa.NPM = request.NPM
	mahasiswa.Nama = request.Nama
	mahasiswa.TanggalLahir = tanggalLahir
	mahasiswa.Fakultas = request.Fakultas
	mahasiswa.Jurusan = request.Jurusan
	mahasiswa.TanggalMasuk = tanggalMasuk
	mahasiswa.JenisKelamin = request.JenisKelamin
	mahasiswa.StatusMahasiswa = request.StatusMahasiswa

	// Update data mahasiswa
	updatedMahasiswa, err := s.mahasiswaRepository.Update(NPM, mahasiswa)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	// Perbarui data user terkait dengan mahasiswa
	user, err := s.userRepository.FindByID(mahasiswa.UserID)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	// Perbarui data user
	user.Username = request.User.Username
	user.Email = request.User.Email

	// Update data user
	updatedUser, err := s.userRepository.Update(user.ID, user)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	// Mengembalikan response mahasiswa yang telah diperbarui beserta data user-nya
	return mapMahasiswaToResponse(updatedMahasiswa, updatedUser), nil
}

// DeleteMahasiswa menghapus data mahasiswa berdasarkan ID
func (s *mahasiswaServiceImpl) DeleteMahasiswa(NPM uint) error {
	// Cari data mahasiswa
	mahasiswa, err := s.mahasiswaRepository.FindByID(NPM)
	if err != nil {
		return utils.ErrNotFound
	}

	// Hapus data mahasiswa
	err = s.mahasiswaRepository.Delete(NPM)
	if err != nil {
		return utils.ErrInternal
	}

	// Hapus data user terkait
	err = s.userRepository.Delete(mahasiswa.User.ID)
	if err != nil {
		return utils.ErrInternal
	}

	return nil
}

// mapMahasiswaToResponse memetakan entitas Mahasiswa dan User ke DTO MahasiswaResponse
func mapMahasiswaToResponse(mahasiswa entities.Mahasiswa, user entities.User) entities.MahasiswaResponse {
	return entities.MahasiswaResponse{
		NPM:             mahasiswa.NPM,
		Nama:            mahasiswa.Nama,
		TanggalLahir:    mahasiswa.TanggalLahir.Format("2006-01-02"),
		Fakultas:        mahasiswa.Fakultas,
		Jurusan:         mahasiswa.Jurusan,
		TanggalMasuk:    mahasiswa.TanggalMasuk.Format("2006-01-02"),
		JenisKelamin:    mahasiswa.JenisKelamin,
		StatusMahasiswa: mahasiswa.StatusMahasiswa,
		User: entities.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}
}
