package entities

import (
	"time"
)

type PenyewaKamar struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	NPM           uint       `gorm:"not null" json:"npm"`
	KamarID       uint       `gorm:"not null" json:"kamar_id"`
	TanggalMulai  time.Time  `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalKeluar *time.Time `gorm:"type:date" json:"tanggal_keluar"`
	Status        string     `gorm:"type:text;not null" json:"status"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	Mahasiswa Mahasiswa `gorm:"foreignKey:NPM" json:"mahasiswa"`
	Kamar     Kamar     `gorm:"foreignKey:KamarID" json:"kamar"`
}

func (PenyewaKamar) TableName() string {
	return "penyewa_kamar"
}
