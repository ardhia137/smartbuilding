package services

import "smartbuilding/entities"

type AuthService interface {
	Login(email, password string) (entities.LoginResponse, error)
	ValidateToken(token string) (*entities.User, error)
	RefreshToken(token string) (entities.LoginResponse, error)
	Logout(token string) error
	ChangePassword(token, old_password, new_password string) error
}
