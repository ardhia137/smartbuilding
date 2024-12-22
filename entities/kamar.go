package entities

import (
	"time"
)

type Kamar struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	NoKamar   uint      `gorm:"not null"`
	Lantai    uint      `gorm:"not null;"`
	Kapasitas uint      `gorm:"not null;"`
	Status    string    `gorm:"type:enum('tersedia', 'tidak tersedia');not null"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Kamar) TableName() string {
	return "kamar"
}
