package services

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
	"strconv"
	"strings"
)

type authServiceImpl struct {
	authRepo    repositories.AuthRepository
	settingRepo repositories.SettingRepository
}

func NewAuthService(authRepo repositories.AuthRepository, settingRepo repositories.SettingRepository) services.AuthService {
	return &authServiceImpl{authRepo, settingRepo}
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

	var setting []entities.Setting
	if user.Role == "admin" {
		setting, err = s.settingRepo.FindAll()
	} else {
		setting, err = s.settingRepo.FindByUserId(user.ID)
	}
	if err != nil {
		return entities.LoginResponse{}, errors.New("setting not found")
	}

	return entities.LoginResponse{
		Token:   token,
		Role:    user.Role,
		UserId:  strconv.FormatUint(uint64(user.ID), 10),
		Setting: setting,
	}, nil
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
func (s *authServiceImpl) ChangePassword(token, oldPassword, newPassword string) error {
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := utils.VerifyToken(token)
	fmt.Println(err)
	if err != nil {
		return errors.New("invalid token")
	}

	user, err := s.authRepo.FindUserByEmail(claims.Email)
	if err != nil {
		return errors.New("user not found")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("old password is incorrect")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	user.Password = string(hashedPassword)
	err = s.authRepo.ChangePassword(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
