package usecases

import (
	"smartbuilding/entities"
)

type AuthUseCase interface {
	Login(email, password string) (entities.LoginResponse, error)
	ValidateToken(token string) (*entities.User, error)
	RefreshToken(token string) (entities.LoginResponse, error)
	Logout(token string) error
}
