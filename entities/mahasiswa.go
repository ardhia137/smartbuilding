package entities

import (
	"time"
)

type Mahasiswa struct {
	NPM             uint      `gorm:"primaryKey;autoIncrement" json:"npm"`
	Nama            string    `gorm:"type:varchar(255);not null" json:"nama"`
	TanggalLahir    time.Time `gorm:"type:date;not null" json:"tanggal_lahir"`
	Fakultas        string    `gorm:"type:varchar(255);not null" json:"fakultas"`
	Jurusan         string    `gorm:"type:varchar(255);not null" json:"jurusan"`
	TanggalMasuk    time.Time `gorm:"type:date;not null" json:"tanggal_masuk"`
	JenisKelamin    string    `gorm:"type:enum('perempuan', 'laki-laki');not null" json:"jenis_kelamin"`
	StatusMahasiswa string    `gorm:"type:enum('aktif', 'tidak aktif');not null" json:"status_mahasiswa"`
	UserID          uint      `gorm:"not null" json:"user_id"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (Mahasiswa) TableName() string {
	return "mahasiswa_doc"
}
