package entities

type DataToren struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	MonitoringName string `json:"monitoring_name" gorm:"type:varchar(255)"`
	KapasitasToren int    `json:"kapasitas_toren" gorm:"NOT NULL"`
	IDSetting      int    `gorm:"not null" json:"id_setting"`
	//
	//Setting Setting `gorm:"foreignKey:IDSetting" json:"setting"`
}

func (DataToren) TableName() string {
	return "data_torent"
}
