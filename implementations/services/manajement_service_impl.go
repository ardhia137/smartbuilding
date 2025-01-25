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

type manajementServiceImpl struct {
	manajementRepository repositories.ManajementRepository
	userRepository       repositories.UserRepository
}

func NewManajementService(
	manajementRepo repositories.ManajementRepository,
	userRepo repositories.UserRepository,
) services.ManajementService {
	return &manajementServiceImpl{
		manajementRepository: manajementRepo,
		userRepository:       userRepo,
	}
}

func (s *manajementServiceImpl) GetAllManajement() ([]entities.ManajementResponse, error) {
	manajements, err := s.manajementRepository.FindAll()
	if err != nil {
		return nil, utils.ErrInternal
	}

	var manajementResponses []entities.ManajementResponse
	for _, manajement := range manajements {
		user, err := s.userRepository.FindByID(manajement.UserID)
		if err != nil {
			return nil, utils.ErrInternal
		}
		manajementResponses = append(manajementResponses, mapManajementToResponse(manajement, user))
	}

	return manajementResponses, nil
}

func (s *manajementServiceImpl) GetManajementByID(NIP uint) (entities.ManajementResponse, error) {
	fmt.Print(NIP)
	manajement, err := s.manajementRepository.FindByID(NIP)
	if err != nil {
		return entities.ManajementResponse{}, utils.ErrNotFound
	}

	user, err := s.userRepository.FindByID(manajement.UserID)
	if err != nil {
		return entities.ManajementResponse{}, utils.ErrInternal
	}

	return mapManajementToResponse(manajement, user), nil
}

func (s *manajementServiceImpl) CreateManajement(request entities.CreateManajementRequest) (entities.ManajementResponse, error) {
	var createdManajement entities.Manajement
	var createdUser entities.User

	db := s.manajementRepository.WithTransaction()
	db = db.Debug()

	err := db.Transaction(func(tx *gorm.DB) error {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.User.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := entities.User{
			Username: request.User.Username,
			Email:    request.User.Email,
			Password: string(hashedPassword),
			Role:     "manajement",
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		createdUser = user

		tanggalLahir, err := time.Parse("2006-01-02", request.TanggalLahir)
		if err != nil {
			return err
		}

		manajement := entities.Manajement{
			NIP:          request.NIP,
			Nama:         request.Nama,
			JenisKelamin: request.JenisKelamin,
			TanggalLahir: tanggalLahir,
			UserID:       user.ID,
		}
		if err := tx.Create(&manajement).Error; err != nil {
			return err
		}
		createdManajement = manajement

		return nil
	})

	if err != nil {
		return entities.ManajementResponse{}, utils.ErrInternal
	}

	return mapManajementToResponse(createdManajement, createdUser), nil
}

func (s *manajementServiceImpl) UpdateManajement(NIP uint, request entities.UpdateManajementRequest) (entities.ManajementResponse, error) {
	manajement, err := s.manajementRepository.FindByID(NIP)
	if err != nil {
		return entities.ManajementResponse{}, utils.ErrNotFound
	}

	manajement.NIP = request.NIP
	manajement.Nama = request.Nama
	manajement.JenisKelamin = request.JenisKelamin

	updatedManajement, err := s.manajementRepository.Update(NIP, manajement)
	if err != nil {
		return entities.ManajementResponse{}, utils.ErrInternal
	}

	user, err := s.userRepository.FindByID(manajement.UserID)
	if err != nil {
		return entities.ManajementResponse{}, utils.ErrInternal
	}

	user.Username = request.User.Username
	user.Email = request.User.Email

	updatedUser, err := s.userRepository.Update(user.ID, user)
	if err != nil {
		return entities.ManajementResponse{}, utils.ErrInternal
	}

	return mapManajementToResponse(updatedManajement, updatedUser), nil
}

func (s *manajementServiceImpl) DeleteManajement(NIP uint) error {
	manajement, err := s.manajementRepository.FindByID(NIP)
	if err != nil {
		return utils.ErrNotFound
	}

	err = s.manajementRepository.Delete(NIP)
	if err != nil {
		return utils.ErrInternal
	}

	err = s.userRepository.Delete(manajement.User.ID)
	if err != nil {
		return utils.ErrInternal
	}

	return nil
}

func mapManajementToResponse(manajement entities.Manajement, user entities.User) entities.ManajementResponse {
	return entities.ManajementResponse{
		NIP:          manajement.NIP,
		Nama:         manajement.Nama,
		TanggalLahir: manajement.TanggalLahir.Format("2006-01-02"),
		JenisKelamin: manajement.JenisKelamin,
		User: entities.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
		},
	}
}
