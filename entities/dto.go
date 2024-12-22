package entities

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`                              // Username pengguna
	Email    string `json:"email" binding:"required,email"`                           // Email valid diperlukan
	Password string `json:"password" binding:"required,min=6,max=100"`                // Password dengan panjang minimal 6 karakter
	Role     string `json:"role" binding:"required,oneof=admin manajement mahasiswa"` // Peran harus salah satu dari admin, management, atau mahasiswa
}
type UpdateUserRequest struct {
	Username string `json:"username" binding:"required"`                              // Username pengguna
	Email    string `json:"email" binding:"required,email"`                           // Email valid diperlukan
	Role     string `json:"role" binding:"required,oneof=admin manajement mahasiswa"` // Peran harus salah satu dari admin, management, atau mahasiswa
}

type UserResponse struct {
	ID       uint   `json:"id"`       // ID pengguna
	Username string `json:"username"` // Username pengguna
	Email    string `json:"email"`    // Email pengguna
	Role     string `json:"role"`     // Peran pengguna
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
	User            CreateUserRequest `json:"user"` // Menambahkan data user
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
	User            UpdateUserRequest `json:"user"` // Menambahkan data user
}
type MahasiswaResponse struct {
	NPM             uint         `json:"npm"`              // Nomor Pokok Mahasiswa
	Nama            string       `json:"nama"`             // Nama mahasiswa
	TanggalLahir    string       `json:"tanggal_lahir"`    // Tanggal lahir mahasiswa
	Fakultas        string       `json:"fakultas"`         // Fakultas mahasiswa
	Jurusan         string       `json:"jurusan"`          // Jurusan mahasiswa
	TanggalMasuk    string       `json:"tanggal_masuk"`    // Tanggal masuk mahasiswa
	JenisKelamin    string       `json:"jenis_kelamin"`    // Jenis kelamin mahasiswa
	StatusMahasiswa string       `json:"status_mahasiswa"` // Status mahasiswa
	User            UserResponse `json:"user"`             // Informasi pengguna terkait
}
