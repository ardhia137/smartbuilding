package repositories

import (
	"smartbuilding/entities"
	"smartbuilding/interfaces/repositories"

	"gorm.io/gorm"
)

type TorentRepositoryImpl struct {
	db *gorm.DB
}

func NewTorentRepository(db *gorm.DB) repositories.TorentRepository {
	return &TorentRepositoryImpl{db: db}
}

func (r *TorentRepositoryImpl) Create(torent *entities.Torent) (*entities.Torent, error) {
	if err := r.db.Create(torent).Error; err != nil {
		return nil, err
	}
	return torent, nil
}

func (r *TorentRepositoryImpl) FindAll() ([]entities.Torent, error) {
	var torentList []entities.Torent
	if err := r.db.Find(&torentList).Error; err != nil {
		return nil, err
	}
	return torentList, nil
}

func (r *TorentRepositoryImpl) FindByGedungID(id int) ([]entities.Torent, error) {
	var torentList []entities.Torent
	if err := r.db.Where("id_gedung = ?", id).Find(&torentList).Error; err != nil {
		return nil, err
	}
	return torentList, nil
}

func (r *TorentRepositoryImpl) FindByID(id int) (*entities.Torent, error) {
	var torent entities.Torent
	if err := r.db.First(&torent, id).Error; err != nil {
		return nil, err
	}
	return &torent, nil
}

func (r *TorentRepositoryImpl) Update(torent *entities.Torent) (*entities.Torent, error) {
	if err := r.db.Save(torent).Error; err != nil {
		return nil, err
	}
	return torent, nil
}

func (r *TorentRepositoryImpl) Delete(id int) error {
	if err := r.db.Delete(&entities.Torent{}, id).Error; err != nil {
		return err
	}
	return nil
}
