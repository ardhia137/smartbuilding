package entities

import "time"

type HakAkses struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int       `json:"id_user" gorm:"NOT NULL"`
	GedungID  int       `json:"id_gedung" gorm:"NOT NULL"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (HakAkses) TableName() string {
	return "hak_akses"
}
