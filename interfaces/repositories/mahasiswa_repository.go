package repositories

import (
	"smartbuilding/entities"
)

type MahasiswaRepository interface {
	FindAll() ([]entities.Mahasiswa, error)
	FindByID(id uint) (entities.Mahasiswa, error)
	Create(mahasiswa entities.Mahasiswa) (entities.Mahasiswa, error)
	Update(NPM uint, mahasiswa entities.Mahasiswa) (entities.Mahasiswa, error)
	Delete(NPM uint) error
}
