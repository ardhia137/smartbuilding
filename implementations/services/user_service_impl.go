package services

import (
	"errors"
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userServiceImpl struct {
	userRepository repositories.UserRepository
	hakAksesRepo   repositories.HakAksesRepository
}

func NewUserService(userRepo repositories.UserRepository, hakAkses repositories.HakAksesRepository) services.UserService {
	return &userServiceImpl{userRepo, hakAkses}
}

func (s *userServiceImpl) GetAllUsers(role string, user_id uint) ([]entities.UserResponse, error) {
	users, err := s.userRepository.FindAll(role, user_id)
	if err != nil {
		return nil, utils.ErrInternal
	}
	var userResponses []entities.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, entities.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Password: user.Password,
		})
	}
	return userResponses, nil
}

func (s *userServiceImpl) GetUserByID(id uint) (entities.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return entities.UserResponse{}, utils.ErrNotFound
	}

	return entities.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role, Password: user.Password}, nil
}

func (s *userServiceImpl) CreateFromAdmin(request entities.CreateUserRequest) (entities.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	// Validasi jika role bukan "admin" tapi HakAkses kosong
	if request.Role != "admin" && len(request.HakAkses) == 0 {
		return entities.UserResponse{}, errors.New("hak akses tidak boleh kosong untuk role ini")
	}

	db := s.userRepository.WithTransaction()
	var createdUser entities.User // Deklarasi di luar transaksi

	err = db.Transaction(func(tx *gorm.DB) error {
		user := entities.User{
			Username: request.Username,
			Email:    request.Email,
			Password: string(hashedPassword),
			Role:     request.Role,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		createdUser = user
		fmt.Println(request.HakAkses)

		// Jika role bukan admin, proses HakAkses
		if request.Role != "admin" {
			var hakAksesList []entities.HakAkses
			for _, hakAkses := range request.HakAkses {
				hakAksesList = append(hakAksesList, entities.HakAkses{
					UserId:   int(user.ID), // Gunakan ID yang baru dibuat
					GedungID: hakAkses.GedungID,
				})
			}
			fmt.Println(request.HakAkses)
			if len(hakAksesList) > 0 {
				if err := tx.Create(&hakAksesList).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Role:     createdUser.Role,
	}, nil
}

func (s *userServiceImpl) CreateFromManajement(id uint, request entities.CreateUserRequest) (entities.UserResponse, error) {
	// Hash password
	hakAkses, err := s.hakAksesRepo.FindByUser(int(id))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	// Validasi jika role bukan "admin" tapi HakAkses kosong
	if request.Role == "pengelola" && len(request.HakAkses) > 1 {
		return entities.UserResponse{}, errors.New("Role Ini Hanya Boleh Mengelola 1 Gedung")
	} else if request.Role == "pengelola" && len(request.HakAkses) == 0 {
		return entities.UserResponse{}, errors.New("hak akses tidak boleh kosong untuk role ini")
	}

	db := s.userRepository.WithTransaction()
	var createdUser entities.User

	err = db.Transaction(func(tx *gorm.DB) error {
		user := entities.User{
			Username: request.Username,
			Email:    request.Email,
			Password: string(hashedPassword),
			Role:     request.Role,
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		createdUser = user // Simpan user yang dibuat

		if request.Role == "manajement" {

			var hakAksesList []entities.HakAkses
			for _, hakAksesItem := range hakAkses {
				hakAksesList = append(hakAksesList, entities.HakAkses{
					UserId:   int(user.ID), // Gunakan ID yang baru dibuat
					GedungID: hakAksesItem.GedungID,
				})
			}

			fmt.Println(request.HakAkses)

			if len(hakAksesList) > 0 {
				if err := tx.Create(&hakAksesList).Error; err != nil {
					return err
				}
			}

		} else if request.Role == "pengelola" {
			hakAkses := entities.HakAkses{
				UserId:   int(user.ID),
				GedungID: request.HakAkses[0].GedungID,
			}
			if err := tx.Create(&hakAkses).Error; err != nil {
				return err
			}
		} else {
			return errors.New("Role Tidak Valid")
		}

		return nil
	})

	if err != nil {
		return entities.UserResponse{}, err
	}

	return entities.UserResponse{
		ID:       createdUser.ID,
		Username: createdUser.Username,
		Email:    createdUser.Email,
		Role:     createdUser.Role,
	}, nil
}

func (s *userServiceImpl) UpdateUser(id uint, request entities.CreateUserRequest) (entities.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return entities.UserResponse{}, utils.ErrNotFound
	}

	user.Username = request.Username
	user.Email = request.Email

	if request.Password != user.Password {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return entities.UserResponse{}, utils.ErrInternal
		}
		user.Password = string(hashedPassword)
	}

	user.Role = request.Role

	updatedUser, err := s.userRepository.Update(id, user)
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	return entities.UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Role:     updatedUser.Role,
		Password: updatedUser.Password,
	}, nil
}

func (s *userServiceImpl) DeleteUser(id uint) error {
	_, err := s.userRepository.FindByID(id)
	if err != nil {
		return utils.ErrNotFound
	}

	err = s.userRepository.Delete(id)
	if err != nil {
		return utils.ErrInternal
	}

	return nil
}
