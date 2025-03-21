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
}

type CreateKamarRequest struct {
	ID        uint   `json:"id"`
	NoKamar   uint   `json:"no_kamar" binding:"required"`
	Lantai    uint   `json:"lantai" binding:"required"`
	Kapasitas uint   `json:"kapasitas" binding:"required"`
	Status    string `json:"status" binding:"required,oneof='tersedia' 'tidak tersedia'"`
}

type KamarResponse struct {
	ID        uint   `json:"id"`
	NoKamar   uint   `json:"no_kamar"`
	Lantai    uint   `json:"lantai"`
	Kapasitas uint   `json:"kapasitas"`
	Status    string `json:"status"`
}

type CreateMahasiswaRequest struct {
	NPM             uint              `json:"npm"`
	Nama            string            `json:"nama"`
	TanggalLahir    string            `json:"tanggal_lahir" time_format:"2006-01-02"`
	Fakultas        string            `json:"fakultas"`
	Jurusan         string            `json:"jurusan"`
	TanggalMasuk    string            `json:"tanggal_masuk" time_format:"2006-01-02"`
	JenisKelamin    string            `json:"jenis_kelamin"`
	StatusMahasiswa string            `json:"status_mahasiswa"`
	User            CreateUserRequest `json:"user"`
}

type UpdateMahasiswaRequest struct {
	NPM             uint              `json:"npm"`
	Nama            string            `json:"nama"`
	TanggalLahir    string            `json:"tanggal_lahir" time_format:"2006-01-02"`
	Fakultas        string            `json:"fakultas"`
	Jurusan         string            `json:"jurusan"`
	TanggalMasuk    string            `json:"tanggal_masuk" time_format:"2006-01-02"`
	JenisKelamin    string            `json:"jenis_kelamin"`
	StatusMahasiswa string            `json:"status_mahasiswa"`
	User            UpdateUserRequest `json:"user"`
}

type MahasiswaResponse struct {
	NPM             uint         `json:"npm"`
	Nama            string       `json:"nama"`
	TanggalLahir    string       `json:"tanggal_lahir"`
	Fakultas        string       `json:"fakultas"`
	Jurusan         string       `json:"jurusan"`
	TanggalMasuk    string       `json:"tanggal_masuk"`
	JenisKelamin    string       `json:"jenis_kelamin"`
	StatusMahasiswa string       `json:"status_mahasiswa"`
	User            UserResponse `json:"user"`
}

type CreateManajementRequest struct {
	NIP          uint              `json:"nip"`
	Nama         string            `json:"nama"`
	TanggalLahir string            `json:"tanggal_lahir" time_format:"2006-01-02"`
	JenisKelamin string            `json:"jenis_kelamin"`
	User         CreateUserRequest `json:"user"`
}

type UpdateManajementRequest struct {
	NIP          uint              `json:"nip"`
	Nama         string            `json:"nama"`
	TanggalLahir string            `json:"tanggal_lahir" time_format:"2006-01-02"`
	JenisKelamin string            `json:"jenis_kelamin"`
	User         UpdateUserRequest `json:"user"`
}

type ManajementResponse struct {
	NIP          uint         `json:"nip"`
	Nama         string       `json:"nama"`
	TanggalLahir string       `json:"tanggal_lahir"`
	JenisKelamin string       `json:"jenis_kelamin"`
	User         UserResponse `json:"user"`
}

type CreatePenyewaKamarRequest struct {
	ID           uint   `json:"id"`
	NPM          uint   `json:"npm"`
	KamarID      uint   `json:"kamar_id"`
	TanggalMulai string `json:"tanggal_mulai" time_format:"2006-01-02"`
	Status       string `json:"status"`
}

type UpdatePenyewaKamarRequest struct {
	ID            uint   `json:"id"`
	NPM           uint   `json:"npm"`
	KamarID       uint   `json:"kamar_id"`
	TanggalMulai  string `json:"tanggal_mulai" time_format:"2006-01-02"`
	TanggalKeluar string `json:"tanggal_keluar" time_format:"2006-01-02"`
	Status        string `json:"status"`
}

type PenyewaKamarResponse struct {
	ID            uint              `json:"id"`
	NPM           uint              `json:"npm"`
	KamarID       uint              `json:"kamar_id"`
	TanggalMulai  string            `json:"tanggal_mulai"`
	TanggalKeluar string            `json:"tanggal_keluar"`
	Status        string            `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Mahasiswa     MahasiswaResponse `json:"mahasiswa"`
	Kamar         KamarResponse     `json:"kamar"`
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
