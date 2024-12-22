package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
)

// userServiceImpl adalah struktur implementasi UserService
type userServiceImpl struct {
	userRepository repositories.UserRepository
}

// NewUserService membuat instance dari UserService
func NewUserService(userRepo repositories.UserRepository) services.UserService {
	return &userServiceImpl{userRepo}
}

// GetAllUsers mengembalikan semua pengguna
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

// GetUserByID mengembalikan informasi pengguna berdasarkan ID
func (s *userServiceImpl) GetUserByID(id uint) (entities.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return entities.UserResponse{}, utils.ErrNotFound
	}

	return entities.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Role: user.Role}, nil
}

// CreateUser membuat pengguna baru
func (s *userServiceImpl) CreateUser(request entities.CreateUserRequest) (entities.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	user := entities.User{
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     request.Role,
	}
	createdUser, err := s.userRepository.Create(user)
	fmt.Println("asd" + request.Username)
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	return entities.UserResponse{ID: createdUser.ID, Username: createdUser.Username, Email: createdUser.Email, Role: user.Role}, nil
}

// UpdateUser memperbarui informasi pengguna berdasarkan ID
func (s *userServiceImpl) UpdateUser(id uint, request entities.CreateUserRequest) (entities.UserResponse, error) {
	// Mencari pengguna berdasarkan ID
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return entities.UserResponse{}, utils.ErrNotFound
	}

	// Memperbarui data pengguna
	user.Username = request.Username
	user.Email = request.Email

	// Jika password ada di request, lakukan hashing dan update password
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return entities.UserResponse{}, utils.ErrInternal
		}
		user.Password = string(hashedPassword)
	}

	// Memperbarui role
	user.Role = request.Role

	// Menyimpan perubahan ke database dengan id dan user
	updatedUser, err := s.userRepository.Update(id, user) // Memasukkan id dan user
	if err != nil {
		return entities.UserResponse{}, utils.ErrInternal
	}

	// Mengembalikan user yang telah diperbarui
	return entities.UserResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Role:     updatedUser.Role,
	}, nil
}

// DeleteUser menghapus pengguna berdasarkan ID
func (s *userServiceImpl) DeleteUser(id uint) error {
	// Mencari pengguna berdasarkan ID
	_, err := s.userRepository.FindByID(id)
	if err != nil {
		return utils.ErrNotFound // Jika pengguna tidak ditemukan
	}

	// Menghapus pengguna dari database
	err = s.userRepository.Delete(id)
	if err != nil {
		return utils.ErrInternal // Jika terjadi error saat penghapusan
	}

	return nil // Penghapusan berhasil
}
