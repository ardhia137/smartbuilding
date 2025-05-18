package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) repositories.AuthRepository {
	return &authRepositoryImpl{db}
}

func (r *authRepositoryImpl) FindUserByEmail(email string) (*entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *authRepositoryImpl) ChangePassword(user *entities.User) error {
	if err := r.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
