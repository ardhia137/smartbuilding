package entities

import (
	"time"
)

type MonitoringData struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	MonitoringName  string    `gorm:"type:varchar(255);not null" json:"monitoring_name"`
	MonitoringValue string    `gorm:"type:varchar(50);not null" json:"monitoring_value"`
	IDSetting       uint      `gorm:"not null" json:"id_setting"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`

	Setting Setting `gorm:"foreignKey:IDSetting" json:"setting"`
}
