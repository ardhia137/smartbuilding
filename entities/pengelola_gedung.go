package entities

type PengelolaGedung struct {
	ID        uint `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int  `json:"id_user" gorm:"NOT NULL"`
	SettingID int  `json:"id_setting" gorm:"NOT NULL"`
}

func (PengelolaGedung) TableName() string {
	return "pengelola_gedung"
}
