package repositories

import (
	"gorm.io/gorm"
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
	return &userRepositoryImpl{db}
}

func (r *userRepositoryImpl) FindAll() ([]entities.User, error) {
	var users []entities.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepositoryImpl) FindByID(id uint) (entities.User, error) {
	var user entities.User
	err := r.db.First(&user, id).Error
	return user, err
}

func (r *userRepositoryImpl) Create(user entities.User) (entities.User, error) {
	err := r.db.Create(&user).Error
	return user, err
}

func (r *userRepositoryImpl) Update(id uint, user entities.User) (entities.User, error) {
	var existingUser entities.User
	err := r.db.First(&existingUser, id).Error
	if err != nil {
		return entities.User{}, err
	}
	user.ID = existingUser.ID
	err = r.db.Save(&user).Error
	return user, err
}

func (r *userRepositoryImpl) Delete(id uint) error {
	var user entities.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return err
	}
	err = r.db.Delete(&user).Error
	return err
}
