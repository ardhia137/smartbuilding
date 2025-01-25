package entities

type Setting struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	HaosURL   string `json:"haos_url" gorm:"type:varchar(255)"`
	HaosToken string `json:"haos_token" gorm:"type:varchar(255)"`
	Scheduler int    `json:"scheduler" gorm:"type:int"`
}

func (Setting) TableName() string {
	return "setting"
}
