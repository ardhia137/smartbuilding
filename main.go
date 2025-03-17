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
	"smartbuilding/utils"
)

func main() {
	log.Println("Starting application...")

	log.Println("Connecting to the database...")
	config.InitDB()

	log.Println("Initializing User repository, service, and use case...")
	userRepository := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepository)
	userUseCase := usecases.UserUseCase(userService)
	userController := controllers.NewUserController(userUseCase)

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

	log.Println("Initializing penyewa kamar repository, service, and use case...")
	penyewaKamarRepository := repositories.NewPenyewaKamarRepository(config.DB)
	penyewaKamarService := services.NewPenyewaKamarService(penyewaKamarRepository, kamarRepository, userRepository, mahasiswaRepository)
	penyewaKamarUsecase := usecases.PenyewaKamarUseCase(penyewaKamarService)
	penyewaKamarController := controllers.NewPenyewaKamarController(penyewaKamarUsecase)

	log.Println("Initializing auth repository, service, and use case...")
	authRepository := repositories.NewAuthRepository(config.DB)
	authService := services.NewAuthService(authRepository)
	authUsecase := usecases.AuthUseCase(authService)
	authController := controllers.NewAuthController(authUsecase)

	log.Println("Initializing data toren repository, service, and use case...")
	dataTorenRepository := repositories.NewDataTorenRepository(config.DB)
	dataTorenService := services.NewDataTorenService(dataTorenRepository)
	dataTorenUsecase := usecases.DataTorenUseCase(dataTorenService)
	dataTorenController := controllers.NewDataTorenController(dataTorenUsecase)

	log.Println("Initializing setting repository, service, and use case...")
	settingRepository := repositories.NewSettingRepository(config.DB)
	settingService := services.NewSettingService(settingRepository, dataTorenRepository)
	settingUsecase := usecases.SettingUseCase(settingService)
	settingController := controllers.NewSettingController(settingUsecase)

	log.Println("Initializing monitoring data repository, service, and use case...")
	monitoringDataRepository := repositories.NewMonitoringDataRepository(config.DB)
	monitoringDataService := services.NewMonitoringDataService(monitoringDataRepository, dataTorenRepository, settingRepository)
	monitoringDataUsecase := usecases.MonitoringDataUseCase(monitoringDataService)
	monitoringDataController := controllers.NewMonitoringDataController(monitoringDataUsecase)

	log.Println("Starting Monitoring Data cron job in the background...")
	go utils.StartMonitoringDataJob(monitoringDataUsecase, settingUsecase, monitoringDataRepository, settingRepository)

	log.Println("Setting up routes...")
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "300")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
	infrastructure.RegisterUserRoutes(router, userController)
	infrastructure.RegisterKamarRoutes(router, kamarController)
	infrastructure.RegisterMahasiswaRoutes(router, mahasiswaController)
	infrastructure.RegisterManajementRoutes(router, manajementController)
	infrastructure.RegisterPenyewaKamarRoutes(router, penyewaKamarController)
	infrastructure.RegisterAuthRoutes(router, authController)
	infrastructure.RegisterMonitoringDataRoutes(router, monitoringDataController)
	infrastructure.RegisterSettingRoutes(router, settingController)
	infrastructure.RegisterDataTorenRoutes(router, dataTorenController)

	log.Println("Starting server on port 3000...")
	err := router.Run(":3000")

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		log.Println("Server is running on port 3000")
	}
}
