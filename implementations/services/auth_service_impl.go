package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
	"strconv"
)

type authServiceImpl struct {
	authRepo repositories.AuthRepository
}

func NewAuthService(authRepo repositories.AuthRepository) services.AuthService {
	return &authServiceImpl{authRepo}
}

func (s *authServiceImpl) Login(email, password string) (entities.LoginResponse, error) {
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return entities.LoginResponse{}, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return entities.LoginResponse{}, errors.New("invalid password")
	}

	token, err := utils.GenerateToken(user.ID, user.Role, user.Email)
	if err != nil {
		return entities.LoginResponse{}, errors.New("failed to generate token")
	}

	return entities.LoginResponse{Token: token, Role: user.Role, UserId: strconv.Itoa(int(user.ID))}, nil
}

func (s *authServiceImpl) ValidateToken(token string) (*entities.User, error) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	user, err := s.authRepo.FindUserByEmail(claims.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (s *authServiceImpl) RefreshToken(token string) (entities.LoginResponse, error) {
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return entities.LoginResponse{}, errors.New("invalid token")
	}

	newToken, err := utils.GenerateToken(claims.UserID, claims.Role, claims.Email)
	if err != nil {
		return entities.LoginResponse{}, errors.New("failed to generate token")
	}

	return entities.LoginResponse{Token: newToken, Role: claims.Role}, nil
}

func (s *authServiceImpl) Logout(token string) error {
	return nil
}
