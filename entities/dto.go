package entities

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`                              // Username pengguna
	Email    string `json:"email" binding:"required,email"`                           // Email valid diperlukan
	Password string `json:"password" binding:"required,min=6,max=100"`                // Password dengan panjang minimal 6 karakter
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
