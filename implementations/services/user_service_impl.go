package services

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
)

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &userServiceImpl{userRepo}
}

func (s *userServiceImpl) GetAllUsers() ([]entities.UserResponse, error) {
	users, err := s.userRepository.FindAll()
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
		})
	}
	return userResponses, nil
}

func (s *userServiceImpl) GetUserByID(id uint) (entities.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return entities.UserResponse{}, utils.ErrNotFound
	}

	return entities.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role}, nil
}
func (s *userServiceImpl) CreateFromAdmin(request entities.CreateUserRequest) (entities.UserResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	// Validasi jika role bukan "admin" tapi PengelolaGedung kosong
	if request.Role != "admin" && len(request.PengelolaGedung) == 0 {
		return entities.UserResponse{}, errors.New("pengelola gedung tidak boleh kosong untuk role ini")
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

		createdUser = user // Simpan user yang dibuat
		fmt.Println(request.PengelolaGedung)

		// Jika role bukan admin, proses PengelolaGedung
		if request.Role != "admin" {
			var pengelolaGedungList []entities.PengelolaGedung
			for _, pengelolaGedung := range request.PengelolaGedung {
				pengelolaGedungList = append(pengelolaGedungList, entities.PengelolaGedung{
					UserId:    int(user.ID), // Gunakan ID yang baru dibuat
					SettingID: pengelolaGedung.SettingID,
				})
			}
			fmt.Println(request.PengelolaGedung)
			if len(pengelolaGedungList) > 0 {
				if err := tx.Create(&pengelolaGedungList).Error; err != nil {
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

func (s *userServiceImpl) UpdateUser(id uint, request entities.CreateUserRequest) (entities.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return entities.UserResponse{}, utils.ErrNotFound
	}

	user.Username = request.Username
	user.Email = request.Email

	if request.Password != "" {
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
