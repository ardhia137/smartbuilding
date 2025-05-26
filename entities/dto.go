package entities

import (
	"time"
)

type CreateUserRequest struct {
	Username        string                             `json:"username" binding:"required"`
	Email           string                             `json:"email" binding:"required,email"`
	Password        string                             `json:"password" binding:"required,min=6,max=100"`
	Role            string                             `json:"role" binding:"required,oneof=admin manajement pengelola"`
	PengelolaGedung []CreateUserPengelolaGedungRequest `json:"pengelola_gedung"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin manajement pengelola"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token   string    `json:"token"`
	Role    string    `json:"role"`
	UserId  string    `json:"user_id"`
	Setting []Setting `json:"Setting"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type CreateMonitoringDataRequest struct {
	MonitoringName  string `json:"monitoring_name" validate:"required"`
	MonitoringValue string `json:"monitoring_value" validate:"required"`
	IDSetting       uint   `json:"id_setting" validate:"required"`
}

type MonitoringDataResponse struct {
	ID              uint      `json:"id"`
	MonitoringName  string    `json:"monitoring_name"`
	MonitoringValue string    `json:"monitoring_value"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type KapasitasTorenData struct {
	Nama           string    `json:"nama"`
	Kapasitas      string    `json:"kapasitas"`
	KapasitasToren string    `json:"kapasitas_toren"`
	VolumeSensor   string    `json:"volume_sensor"`
	CreatedAt      time.Time `json:"created_at"`
}

type GetAirDataResponse struct {
	NamaGedung             string                     `json:"nama_gedung"`
	KapasitasToren         []KapasitasTorenData       `json:"kapasitasToren"`
	AirKeluar              string                     `json:"AirKeluar"`
	AirMasuk               string                     `json:"AirMasuk"`
	DataPenggunaanHarian   map[string][]PenggunaanAir `json:"DataPenggunaanHarian"`
	DataPenggunaanMingguan map[string][]PenggunaanAir `json:"DataPenggunaanMingguan"`
	DataPenggunaanBulanan  map[string][]PenggunaanAir `json:"DataPenggunaanBulanan"`
	DataPenggunaanTahunan  map[string][]PenggunaanAir `json:"DataPenggunaanTahunan"`
	CreatedAt              time.Time                  `json:"CreatedAt"`
	UpdatedAt              time.Time                  `json:"UpdatedAt"`
}

type PenggunaanAir struct {
	Pipa   string `json:"pipa"`
	Volume string `json:"volume"`
}

type GetListrikDataResponse struct {
	NamaGedung                    string                         `json:"nama_gedung"`
	TotalWatt                     string                         `json:"TotalWatt"`
	TotalDayaListrik              []TotalDayaListrik             `json:"TotalDayaListrik"`
	BiayaPemakaian                []BiayaListrik                 `json:"BiayaPemakaian"`
	DataPenggunaanListrikHarian   map[string][]PenggunaanListrik `json:"DataPenggunaanListrikHarian"`
	DataBiayaListrikHarian        map[string][]BiayaListrik      `json:"DataBiayaListrikHarian"`
	DataPenggunaanListrikMingguan map[string][]PenggunaanListrik `json:"DataPenggunaanListrikMingguan"`
	DataBiayaListrikMingguan      map[string][]BiayaListrik      `json:"DataBiayaListrikMingguan"`
	DataPenggunaanListrikBulanan  map[string][]PenggunaanListrik `json:"DataPenggunaanListrikBulanan"`
	DataBiayaListrikBulanan       map[string][]BiayaListrik      `json:"DataBiayaListrikBulanan"`
	DataPenggunaanListrikTahunan  map[string][]PenggunaanListrik `json:"DataPenggunaanListrikTahunan"`
	DataBiayaListrikTahunan       map[string][]BiayaListrik      `json:"DataBiayaListrikTahunan"`
	CreatedAt                     time.Time                      `json:"CreatedAt"`
	UpdatedAt                     time.Time                      `json:"UpdatedAt"`
}

type PenggunaanListrik struct {
	Nama  string `json:"nama"`
	Value string `json:"Value"`
}
type TotalDayaListrik struct {
	Nama  string `json:"nama"`
	Value string `json:"Value"`
}

type BiayaListrik struct {
	Nama  string `json:"Nama"`
	Biaya string `json:"Biaya"`
}

type CreateSettingRequest struct {
	NamaGedung   string                   `json:"nama_gedung" binding:"required"`
	HaosURL      string                   `json:"haos_url"       binding:"required"`
	HaosToken    string                   `json:"haos_token"     binding:"required"`
	Scheduler    int                      `json:"scheduler"      binding:"required"`
	HargaListrik int                      `json:"harga_listrik"  binding:"required"`
	DataToren    []CreateDataTorenRequest `json:"data_toren"`
	JenisListrik string                   `json:"jenis_listrik" binding:"required,oneof='1_phase' '3_phase'"`
}

type UpdateSettingRequest struct {
	NamaGedung   string `json:"nama_gedung" binding:"required"`
	HaosURL      string `json:"haos_url"       binding:"required"`
	HaosToken    string `json:"haos_token"     binding:"required"`
	Scheduler    int    `json:"scheduler"      binding:"required"`
	HargaListrik int    `json:"harga_listrik"  binding:"required"`
	JenisListrik string `json:"jenis_listrik" binding:"required,oneof='1_phase' '3_phase'"`
}

type SettingResponse struct {
	ID           int    `json:"id"`
	NamaGedung   string `json:"nama_gedung"`
	HaosURL      string `json:"haos_url"`
	HaosToken    string `json:"haos_token"`
	Scheduler    int    `json:"scheduler"`
	HargaListrik int    `json:"harga_listrik"`
	JenisListrik string `json:"jenis_listrik" binding:"required"`
}

type SettingResponseCreate struct {
	ID           int         `json:"id"`
	NamaGedung   string      `json:"nama_gedung"`
	HaosURL      string      `json:"haos_url"`
	HaosToken    string      `json:"haos_token"`
	Scheduler    int         `json:"scheduler"`
	HargaListrik int         `json:"harga_listrik"`
	JenisListrik string      `json:"jenis_listrik" binding:"required"`
	DataToren    []DataToren `json:"data_toren"`
}

type CreateDataTorenRequest struct {
	MonitoringName string `json:"monitoring_name" binding:"required"`
	KapasitasToren int    `json:"kapasitas_toren" binding:"required"`
	IDSetting      int    `json:"id_setting" binding:"required"`
}

type DataTorenResponse struct {
	ID             uint   `json:"id"`
	MonitoringName string `json:"monitoring_name"`
	KapasitasToren int    `json:"kapasitas_toren"`
	IDSetting      int    `json:"id_setting"`
}

type CreatePengelolaGedungRequest struct {
	UserID    int `json:"id_user" binding:"required"`
	SettingID int `json:"id_setting" binding:"required"`
}

type CreateUserPengelolaGedungRequest struct {
	SettingID int `json:"setting_id" binding:"required"`
}

type PengelolaGedungResponse struct {
	ID        uint `json:"id"`
	UserID    int  `json:"id_user"`
	SettingID int  `json:"id_setting"`
}

type AllPengelolaGedungResponse struct {
	ID         uint   `json:"id"`
	NamaGedung string `json:"nama_gedung"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	SettingID  int    `json:"setting_id"`
}
