package usecases

import "smartbuilding/entities"

type UserUseCase interface {
	GetAllUsers() ([]entities.UserResponse, error)
	GetUserByID(id uint) (entities.UserResponse, error)
	CreateFromAdmin(request entities.CreateUserRequest) (entities.UserResponse, error)
	UpdateUser(id uint, request entities.CreateUserRequest) (entities.UserResponse, error)
	DeleteUser(id uint) error
}
