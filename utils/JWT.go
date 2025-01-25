package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var jwtKey = []byte("Open404!")

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, role string, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
