package entities

type Setting struct {
	ID           int    `json:"id" gorm:"primaryKey;autoIncrement"`
	NamaGedung   string `json:"nama_gedung" gorm:"type:varchar(255)"`
	HaosURL      string `json:"haos_url" gorm:"type:varchar(255)"`
	HaosToken    string `json:"haos_token" gorm:"type:varchar(255)"`
	Scheduler    int    `json:"scheduler" gorm:"type:int"`
	HargaListrik int    `json:"harga_listrik" gorm:"type:int"`
	JenisListrik string `json:"jenis_listrik" binding:"required,oneof='1_phase' '3_phase'"`
}

func (Setting) TableName() string {
	return "setting"
}
