package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"smartbuilding/config"
	"smartbuilding/controllers"
	"smartbuilding/implementations/repositories"
	"smartbuilding/implementations/services"
	"smartbuilding/infrasturcture"
	"smartbuilding/usecases"
)

func main() {
	// Logging untuk memulai aplikasi
	log.Println("Starting application...")

	// Inisialisasi koneksi database
	log.Println("Connecting to the database...")
	config.InitDB()

	// Inisialisasi Repository, Service, dan UseCase untuk User
	log.Println("Initializing User repository, service, and use case...")
	userRepository := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepository)
	userUseCase := usecases.UserUseCase(userService)
	userController := controllers.NewUserController(userUseCase)

	// Inisialisasi Repository, Service, dan UseCase untuk User
	log.Println("Initializing User repository, service, and use case...")
	kamarRepository := repositories.NewKamarRepository(config.DB)
	kamarService := services.NewKamarService(kamarRepository)
	kamarUsecase := usecases.KamarUseCase(kamarService)
	kamarController := controllers.NewKamarController(kamarUsecase)

	log.Println("Initializing User repository, service, and use case...")
	mahasiswaRepository := repositories.NewMahasiswaRepository(config.DB)
	mahasiswaService := services.NewMahasiswaService(mahasiswaRepository, userRepository)
	mahasiswaUsecase := usecases.MahasiswaUseCase(mahasiswaService)
	mahasiswaController := controllers.NewMahasiswaController(mahasiswaUsecase)

	log.Println("Initializing User repository, service, and use case...")
	manajementRepository := repositories.NewManajementRepository(config.DB)
	manajementService := services.NewManajementService(manajementRepository, userRepository)
	manajementUsecase := usecases.ManajementUseCase(manajementService)
	manajementController := controllers.NewManajementController(manajementUsecase)

	// Set up Router
	log.Println("Setting up routes...")
	router := gin.Default()
	infrastructure.RegisterUserRoutes(router, userController)
	infrastructure.RegisterKamarRoutes(router, kamarController)
	infrastructure.RegisterMahasiswaRoutes(router, mahasiswaController)
	infrastructure.RegisterManajementRoutes(router, manajementController)

	// Jalankan server pada port 3000
	log.Println("Starting server on port 3000...")
	err := router.Run(":3000")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		log.Println("Server is running on port 3000")
	}

}
