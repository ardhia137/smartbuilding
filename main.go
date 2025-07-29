package main

import (
	"log"
	"smartbuilding/config"
	"smartbuilding/controllers"
	"smartbuilding/implementations/repositories"
	"smartbuilding/implementations/services"
	infrastructure "smartbuilding/infrasturcture"
	"smartbuilding/usecases"
	"smartbuilding/utils"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Starting application...")

	log.Println("Connecting to the database...")
	config.InitDB()

	hakAksesRepository := repositories.NewHakAksesRepository(config.DB)
	hakAksesService := services.NewHakAksesService(hakAksesRepository)
	hakAksesUsecase := usecases.HakAksesUseCase(hakAksesService)
	hakAksesController := controllers.NewHakAksesController(hakAksesUsecase)

	log.Println("Initializing User repository, service, and use case...")
	userRepository := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepository, hakAksesRepository)
	userUseCase := usecases.UserUseCase(userService)
	userController := controllers.NewUserController(userUseCase)

	//log.Println("Initializing data toren repository, service, and use case...")
	torentRepository := repositories.NewTorentRepository(config.DB)
	torentService := services.NewTorentService(torentRepository)
	torentUsecase := usecases.TorentUseCase(torentService)
	torentController := controllers.NewTorentController(torentUsecase)

	log.Println("Initializing gedung repository, service, and use case...")
	gedungRepository := repositories.NewGedungRepository(config.DB)
	gedungService := services.NewGedungService(gedungRepository, torentRepository)
	gedungUsecase := usecases.GedungUseCase(gedungService)
	gedungController := controllers.NewGedungController(gedungUsecase)

	log.Println("Initializing monitoring log repository, service, and use case...")
	monitoringLogRepository := repositories.NewMonitoringLogRepository(config.DB)
	monitoringLogService := services.NewMonitoringLogService(monitoringLogRepository, torentRepository, gedungRepository)
	monitoringLogUsecase := usecases.MonitoringLogUseCase(monitoringLogService)
	monitoringLogController := controllers.NewMonitoringLogController(monitoringLogUsecase, hakAksesUsecase)

	log.Println("Initializing auth repository, service, and use case...")
	authRepository := repositories.NewAuthRepository(config.DB)
	authService := services.NewAuthService(authRepository, gedungRepository)
	authUsecase := usecases.AuthUseCase(authService)
	authController := controllers.NewAuthController(authUsecase)

	log.Println("Starting Monitoring Log cron job in the background...")
	go utils.StartMonitoringLogJob(monitoringLogUsecase, gedungUsecase, monitoringLogRepository, gedungRepository)

	log.Println("Setting up routes...")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Tambahkan ini supaya OPTIONS tidak 404
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})
	infrastructure.RegisterUserRoutes(router, userController)
	infrastructure.RegisterAuthRoutes(router, authController)
	infrastructure.RegisterMonitoringLogRoutes(router, monitoringLogController)
	infrastructure.RegisterGedungRoutes(router, gedungController)
	infrastructure.RegisterTorentRoutes(router, torentController)
	infrastructure.RegisterHakAksesRoutes(router, hakAksesController)

	log.Println("Starting server on port 1312...")
	err := router.Run(":1312")

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	} else {
		log.Println("Server is running on port 3000")
	}
}
