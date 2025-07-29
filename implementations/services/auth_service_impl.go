package services

import (
	"errors"
	"fmt"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
	"smartbuilding/interfaces/services"
	"smartbuilding/utils"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type authServiceImpl struct {
	authRepo          repositories.AuthRepository
	gedungRepo        repositories.GedungRepository
	blacklistedTokens map[string]time.Time
	blacklistMutex    sync.RWMutex
}

func NewAuthService(authRepo repositories.AuthRepository, gedungRepo repositories.GedungRepository) services.AuthService {
	service := &authServiceImpl{
		authRepo:          authRepo,
		gedungRepo:        gedungRepo,
		blacklistedTokens: make(map[string]time.Time),
	}

	// Jalankan cleanup routine setiap jam
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			service.cleanupExpiredBlacklistedTokens()
		}
	}()

	return service
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

	var gedung []entities.Gedung
	if user.Role == "admin" {
		gedung, err = s.gedungRepo.FindAll()
	} else {
		gedung, err = s.gedungRepo.FindByUserId(user.ID)
	}
	if err != nil {
		return entities.LoginResponse{}, errors.New("gedung not found")
	}

	return entities.LoginResponse{
		Token:  token,
		Role:   user.Role,
		UserId: strconv.FormatUint(uint64(user.ID), 10),
		Gedung: gedung,
	}, nil
}

func (s *authServiceImpl) ValidateToken(token string) (*entities.User, error) {
	// Bersihkan token dari "Bearer " prefix
	token = strings.TrimPrefix(token, "Bearer ")

	// Cek apakah token ada di blacklist
	if s.isTokenBlacklisted(token) {
		return nil, errors.New("token has been revoked (user logged out)")
	}

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
	token = strings.TrimPrefix(token, "Bearer ")
	claims, err := utils.VerifyToken(token)
	if err != nil {
		return errors.New("invalid token")
	}
	s.blacklistMutex.Lock()
	s.blacklistedTokens[token] = time.Now()
	s.blacklistMutex.Unlock()
	fmt.Printf("Token for user %s (ID: %d) has been blacklisted\n", claims.Email, claims.UserID)

	return nil
}

func (s *authServiceImpl) isTokenBlacklisted(token string) bool {
	s.blacklistMutex.RLock()
	defer s.blacklistMutex.RUnlock()

	_, exists := s.blacklistedTokens[token]
	return exists
}

func (s *authServiceImpl) cleanupExpiredBlacklistedTokens() {
	s.blacklistMutex.Lock()
	defer s.blacklistMutex.Unlock()

	now := time.Now()
	for token, logoutTime := range s.blacklistedTokens {
		if now.Sub(logoutTime) > 24*time.Hour {
			delete(s.blacklistedTokens, token)
		}
	}
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
