package entities

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Username  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);not null;unique"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      string    `gorm:"type:enum('admin', 'manajement', 'mahasiswa');default:'';not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// Method TableName untuk menetapkan nama tabel yang berbeda
func (User) TableName() string {
	return "user" // Nama tabel yang Anda inginkan
}
