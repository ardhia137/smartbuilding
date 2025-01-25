package entities

import (
	"time"
)

type MonitoringDataHarian struct {
	ID              uint      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	MonitoringName  string    `gorm:"type:varchar(255);not null" json:"monitoring_name"`
	MonitoringValue string    `gorm:"type:varchar(50);not null" json:"monitoring_value"`
	CreatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}
