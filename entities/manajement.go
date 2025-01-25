package entities

import (
	"time"
)

type Manajement struct {
	NIP          uint      `gorm:"primaryKey;autoIncrement;column:nip" json:"nip"`
	Nama         string    `gorm:"type:varchar(255);not null" json:"nama"`
	TanggalLahir time.Time `gorm:"type:date;not null" json:"tanggal_lahir"`
	JenisKelamin string    `gorm:"type:enum('perempuan', 'laki-laki');not null" json:"jenis_kelamin"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`

	User User `gorm:"foreignKey:UserID" json:"user"`
}

func (Manajement) TableName() string {
	return "manajement_doc"
}
