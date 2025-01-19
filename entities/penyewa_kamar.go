package entities

import (
	"time"
)

type PenyewaKamar struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	NPM           uint       `gorm:"not null" json:"npm"`      // Relasi ke tabel kamar
	KamarID       uint       `gorm:"not null" json:"kamar_id"` // Relasi ke tabel user
	TanggalMulai  time.Time  `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalKeluar *time.Time `gorm:"type:date" json:"tanggal_keluar"` // Opsional, bisa kosong
	Status        string     `gorm:"type:text;not null" json:"status"`
	CreatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`

	// Relasi
	Mahasiswa Mahasiswa `gorm:"foreignKey:NPM" json:"mahasiswa"` // Asumsi ada entitas User
	Kamar     Kamar     `gorm:"foreignKey:KamarID" json:"kamar"`
}

func (PenyewaKamar) TableName() string {
	return "penyewa_kamar" // Sesuaikan dengan nama tabel yang digunakan
}
