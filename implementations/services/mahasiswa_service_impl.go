package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
	"time"
)

type mahasiswaServiceImpl struct {
	mahasiswaRepository repositories.MahasiswaRepository
	userRepository      repositories.UserRepository
}

func NewMahasiswaService(
	mahasiswaRepo repositories.MahasiswaRepository,
	userRepo repositories.UserRepository,
) services.MahasiswaService {
	return &mahasiswaServiceImpl{
		mahasiswaRepository: mahasiswaRepo,
		userRepository:      userRepo,
	}
}

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

func (s *mahasiswaServiceImpl) CreateMahasiswa(request entities.CreateMahasiswaRequest) (entities.MahasiswaResponse, error) {
	var createdMahasiswa entities.Mahasiswa
	var createdUser entities.User

	db := s.mahasiswaRepository.WithTransaction()

	err := db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.User.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := entities.User{
			Username: request.User.Username,
			Email:    request.User.Email,
			Password: string(hashedPassword),
			Role:     "mahasiswa",
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		createdUser = user

		tanggalMasuk, err := time.Parse("2006-01-02", request.TanggalMasuk)
		if err != nil {
			return err
		}
		tanggalLahir, err := time.Parse("2006-01-02", request.TanggalLahir)
		if err != nil {
			return err
		}

		mahasiswa := entities.Mahasiswa{
			NPM:             request.NPM,
			Nama:            request.Nama,
			TanggalLahir:    tanggalLahir,
			Fakultas:        request.Fakultas,
			Jurusan:         request.Jurusan,
			TanggalMasuk:    tanggalMasuk,
			JenisKelamin:    request.JenisKelamin,
			StatusMahasiswa: request.StatusMahasiswa,
			UserID:          user.ID,
		}
		if err := tx.Create(&mahasiswa).Error; err != nil {
			return err
		}
		createdMahasiswa = mahasiswa

		return nil
	})

	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	return mapMahasiswaToResponse(createdMahasiswa, createdUser), nil
}

func (s *mahasiswaServiceImpl) UpdateMahasiswa(NPM uint, request entities.UpdateMahasiswaRequest) (entities.MahasiswaResponse, error) {
	mahasiswa, err := s.mahasiswaRepository.FindByID(NPM)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrNotFound
	}
	tanggalMasuk, err := time.Parse("2006-01-02", request.TanggalMasuk)
	tanggalLahir, err := time.Parse("2006-01-02", request.TanggalMasuk)
	fmt.Print(request.TanggalMasuk)
	fmt.Print(tanggalMasuk)

	mahasiswa.NPM = request.NPM
	mahasiswa.Nama = request.Nama
	mahasiswa.TanggalLahir = tanggalLahir
	mahasiswa.Fakultas = request.Fakultas
	mahasiswa.Jurusan = request.Jurusan
	mahasiswa.TanggalMasuk = tanggalMasuk
	mahasiswa.JenisKelamin = request.JenisKelamin
	mahasiswa.StatusMahasiswa = request.StatusMahasiswa

	updatedMahasiswa, err := s.mahasiswaRepository.Update(NPM, mahasiswa)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	user, err := s.userRepository.FindByID(mahasiswa.UserID)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	user.Username = request.User.Username
	user.Email = request.User.Email

	updatedUser, err := s.userRepository.Update(user.ID, user)
	if err != nil {
		return entities.MahasiswaResponse{}, utils.ErrInternal
	}

	return mapMahasiswaToResponse(updatedMahasiswa, updatedUser), nil
}

func (s *mahasiswaServiceImpl) DeleteMahasiswa(NPM uint) error {
	mahasiswa, err := s.mahasiswaRepository.FindByID(NPM)
	if err != nil {
		return utils.ErrNotFound
	}

	err = s.mahasiswaRepository.Delete(NPM)
	if err != nil {
		return utils.ErrInternal
	}

	err = s.userRepository.Delete(mahasiswa.User.ID)
	if err != nil {
		return utils.ErrInternal
	}

	return nil
}

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
