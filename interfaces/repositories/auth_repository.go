package repositories

import (
	"smartbuilding/entities"
)

type AuthRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
}
