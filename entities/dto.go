package entities

import (
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Role     string `json:"role" binding:"required,oneof=admin manajement mahasiswa"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin manajement mahasiswa"`
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
	Token string `json:"token"`
	Role  string `json:"role"`
}

type CreateMonitoringDataRequest struct {
	MonitoringName  string `json:"monitoring_name" validate:"required"`
	MonitoringValue string `json:"monitoring_value" validate:"required"`
}

type MonitoringDataResponse struct {
	ID              uint      `json:"id"`
	MonitoringName  string    `json:"monitoring_name"`
	MonitoringValue string    `json:"monitoring_value"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GetAirDataResponse struct {
	KapasitasToren         string                     `json:"KapasitasToren"`
	AirKeluar              string                     `json:"AirKeluar"`
	AirMasuk               string                     `json:"AirMasuk"`
	VolumeSensor           string                     `json:"VolumeSensor"`
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
	TotalWatt                     string                         `json:"TotalWatt"`
	TotalDayaListrikLT1           string                         `json:"TotalDayaListrikLT1"`
	TotalDayaListrikLT2           string                         `json:"TotalDayaListrikLT2"`
	TotalDayaListrikLT3           string                         `json:"TotalDayaListrikLT3"`
	TotalDayaListrikLT4           string                         `json:"TotalDayaListrikLT4"`
	BiayaPemakaianLT1             string                         `json:"BiayaPemakaianLT1"`
	BiayaPemakaianLT2             string                         `json:"BiayaPemakaianLT2"`
	BiayaPemakaianLT3             string                         `json:"BiayaPemakaianLT3"`
	BiayaPemakaianLT4             string                         `json:"BiayaPemakaianLT4"`
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
	Lantai int    `json:"Lantai"`
	Value  string `json:"Value"`
}

type BiayaListrik struct {
	Lantai int    `json:"Lantai"`
	Biaya  string `json:"Biaya"`
}

type CreateSettingRequest struct {
	HaosURL   string `json:"haos_url" binding:"required"`
	HaosToken string `json:"haos_token" binding:"required"`
	Scheduler int    `json:"scheduler" binding:"required"`
}

type SettingResponse struct {
	ID        int    `json:"id"`
	HaosURL   string `json:"haos_url"`
	HaosToken string `json:"haos_token"`
	Scheduler int    `json:"scheduler"`
}
