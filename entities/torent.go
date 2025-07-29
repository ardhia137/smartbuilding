package entities

import "time"

type Torent struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	MonitoringName string    `json:"monitoring_name" gorm:"type:varchar(255)"`
	KapasitasToren int       `json:"kapasitas_toren" gorm:"NOT NULL"`
	IDGedung       int       `gorm:"not null" json:"id_gedung"`
	CreatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	//Gedung Gedung `gorm:"foreignKey:IDGedung" json:"gedung"`
}

func (Torent) TableName() string {
	return "torent"
}
