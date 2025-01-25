package services

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
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
