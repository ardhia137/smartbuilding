package services

import (
	"smartbuilding/entities"
)

type UserService interface {
	GetAllUsers() ([]entities.UserResponse, error)
	GetUserByID(id uint) (entities.UserResponse, error)
	CreateUser(request entities.CreateUserRequest) (entities.UserResponse, error)
	UpdateUser(id uint, request entities.CreateUserRequest) (entities.UserResponse, error)
	DeleteUser(id uint) error
}
