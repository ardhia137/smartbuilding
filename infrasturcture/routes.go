package infrastructure

import (
	"github.com/gin-gonic/gin"
	"smartbuilding/controllers"
	"smartbuilding/utils"
)

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	apiGroup := router.Group("/api")
	userRoutes := apiGroup.Group("/users")
	userRoutes.Use(utils.RoleMiddleware("admin"))
	{
		userRoutes.GET("", userController.GetAllUsers)
		userRoutes.POST("", userController.CreateUser)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
}

func RegisterKamarRoutes(router *gin.Engine, kamarController *controllers.KamarController) {
	apiGroup := router.Group("/api")
	kamarRoutes := apiGroup.Group("/kamar")
	kamarRoutes.Use(utils.RoleMiddleware("admin", "manajement"))
	{
		kamarRoutes.GET("", kamarController.GetAllKamars)
		kamarRoutes.POST("", kamarController.CreateKamar)
		kamarRoutes.GET("/:id", kamarController.GetKamarByID)
		kamarRoutes.PUT("/:id", kamarController.UpdateKamar)
		kamarRoutes.DELETE("/:id", kamarController.DeleteKamar)
	}
}

func RegisterMahasiswaRoutes(router *gin.Engine, mahasiswaController *controllers.MahasiswaController) {
	apiGroup := router.Group("/api")
	mahasiswaRoutes := apiGroup.Group("/mahasiswa")
	mahasiswaRoutes.Use(utils.RoleMiddleware("admin", "manajement"))
	{
		mahasiswaRoutes.GET("", mahasiswaController.GetAllMahasiswas)
		mahasiswaRoutes.POST("", mahasiswaController.CreateMahasiswa)
		mahasiswaRoutes.GET("/:NPM", mahasiswaController.GetMahasiswaByID)
		mahasiswaRoutes.PUT("/:NPM", mahasiswaController.UpdateMahasiswa)
		mahasiswaRoutes.DELETE("/:NPM", mahasiswaController.DeleteMahasiswa)
	}
}

func RegisterManajementRoutes(router *gin.Engine, manajementController *controllers.ManajementController) {
	apiGroup := router.Group("/api")
	manajementRoutes := apiGroup.Group("/manajement")
	manajementRoutes.Use(utils.RoleMiddleware("admin", "manajement"))
	{
		manajementRoutes.GET("", manajementController.GetAllManajements)
		manajementRoutes.POST("", manajementController.CreateManajement)
		manajementRoutes.GET("/:NIP", manajementController.GetManajementByID)
		manajementRoutes.PUT("/:NIP", manajementController.UpdateManajement)
		manajementRoutes.DELETE("/:NIP", manajementController.DeleteManajement)
	}
}

func RegisterPenyewaKamarRoutes(router *gin.Engine, penyewaKamarController *controllers.PenyewaKamarController) {
	apiGroup := router.Group("/api")
	penyewaKamarRoutes := apiGroup.Group("/penyewa_kamar")
	penyewaKamarRoutes.Use(utils.RoleMiddleware("admin", "manajement"))
	{
		penyewaKamarRoutes.GET("", penyewaKamarController.GetAllPenyewaKamars)
		penyewaKamarRoutes.POST("", penyewaKamarController.CreatePenyewaKamar)
		penyewaKamarRoutes.GET("/:id", penyewaKamarController.GetPenyewaKamarByID)
		penyewaKamarRoutes.PUT("/:id", penyewaKamarController.UpdatePenyewaKamar)
		penyewaKamarRoutes.DELETE("/:id", penyewaKamarController.DeletePenyewaKamar)
	}
}
func RegisterAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	apiGroup := router.Group("/api")
	authRoutes := apiGroup.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)
	}
}

func RegisterMonitoringDataRoutes(router *gin.Engine, monitoringDataController *controllers.MonitoringDataController) {
	apiRoutes := router.Group("/api")
	apiRoutes.Use(utils.RoleMiddleware("admin", "manajement"))
	{
		apiRoutes.GET("/monitoring_air", monitoringDataController.GetAirMonitoringData)
		apiRoutes.GET("/monitoring_listrik", monitoringDataController.GetListrikMonitoringData)

	}
}

func RegisterSettingRoutes(router *gin.Engine, settingController *controllers.SettingController) {
	apiGroup := router.Group("/api")
	settingRoutes := apiGroup.Group("/setting")
	settingRoutes.Use(utils.RoleMiddleware("admin"))
	{
		settingRoutes.GET("", settingController.GetAllSetting)
		settingRoutes.POST("", settingController.CreateSetting)
		settingRoutes.GET("/:id", settingController.GetSettingByID)
		settingRoutes.PUT("/:id", settingController.UpdateSetting)
		settingRoutes.DELETE("/:id", settingController.DeleteSetting)

	}
}
