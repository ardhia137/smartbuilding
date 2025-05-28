package infrastructure

import (
	"github.com/gin-gonic/gin"
	"smartbuilding/controllers"
	"smartbuilding/utils"
)

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	apiGroup := router.Group("/api")

	// Routes only for admin and manajement
	adminRoutes := apiGroup.Group("/users")
	adminRoutes.Use(utils.RoleMiddleware("admin", "manajement"), utils.UserIDMiddleware())
	{
		adminRoutes.GET("", userController.GetAllUsers)
		adminRoutes.POST("", userController.CreateUser)
		adminRoutes.GET("/:id", userController.GetUserByID)
		adminRoutes.PUT("/:id", userController.UpdateUser)
		adminRoutes.DELETE("/:id", userController.DeleteUser)
	}

	// Routes for any authenticated user (e.g., pengelola, admin, etc.)
	userSelfRoutes := apiGroup.Group("/users")
	userSelfRoutes.Use(utils.UserIDMiddleware()) // No RoleMiddleware, so any logged-in user can access
	{
		userSelfRoutes.GET("/me", userController.GetMe)
	}
}
func RegisterAuthRoutes(router *gin.Engine, authController *controllers.AuthController) {
	apiGroup := router.Group("/api")
	authRoutes := apiGroup.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/logout", authController.Logout)

	}
	apiGroup.Use(utils.RoleMiddleware("admin", "manajement", "pengelola"), utils.UserIDMiddleware())
	apiGroup.PUT("/change-password", authController.ChangePassword)

}

func RegisterMonitoringDataRoutes(router *gin.Engine, monitoringDataController *controllers.MonitoringDataController) {
	apiRoutes := router.Group("/api")
	apiRoutes.Use(utils.RoleMiddleware("admin", "manajement", "pengelola"), utils.UserIDMiddleware())
	{
		apiRoutes.GET("/monitoring_air/:id", monitoringDataController.GetAirMonitoringData)
		apiRoutes.GET("/monitoring_listrik/:id", monitoringDataController.GetListrikMonitoringData)

	}
}

func RegisterSettingRoutes(router *gin.Engine, settingController *controllers.SettingController) {
	apiGroup := router.Group("/api")
	settingRoutes := apiGroup.Group("/setting")
	settingRoutes.Use(utils.RoleMiddleware("admin", "manajement", "pengelola"), utils.UserIDMiddleware())
	{
		settingRoutes.GET("", settingController.GetAllSetting)
		settingRoutes.POST("", settingController.CreateSetting)
		settingRoutes.GET("/:id", settingController.GetSettingByID)
		settingRoutes.PUT("/:id", settingController.UpdateSetting)
		settingRoutes.DELETE("/:id", settingController.DeleteSetting)
	}
}

func RegisterDataTorenRoutes(router *gin.Engine, dataTorenController *controllers.DataTorenController) {
	apiGroup := router.Group("/api")
	dataTorenRoutes := apiGroup.Group("/data_toren")
	dataTorenRoutes.Use(utils.RoleMiddleware("admin", "manajement"))
	{
		dataTorenRoutes.GET("", dataTorenController.GetAllDataToren)
		dataTorenRoutes.POST("", dataTorenController.CreateDataToren)
		dataTorenRoutes.GET("/:id", dataTorenController.GetDataTorenByID)
		dataTorenRoutes.GET("/setting/:id", dataTorenController.GetDataTorenBySettingID)
		dataTorenRoutes.PUT("/:id", dataTorenController.UpdateDataToren)
		dataTorenRoutes.DELETE("/:id", dataTorenController.DeleteDataToren)

	}
}

func RegisterPengelolaGedungRoutes(router *gin.Engine, pengelolaGedungController *controllers.PengelolaGedungController) {
	apiGroup := router.Group("/api")
	pengelolaGedungRoutes := apiGroup.Group("/pengelola_gedung")
	pengelolaGedungRoutes.Use(utils.RoleMiddleware("admin", "manajement"), utils.UserIDMiddleware())
	{
		pengelolaGedungRoutes.GET("", pengelolaGedungController.GetAllPengelolaGedung)
		pengelolaGedungRoutes.POST("", pengelolaGedungController.CreatePengelolaGedung)
		pengelolaGedungRoutes.GET("/:id", pengelolaGedungController.GetPengelolaGedungBySettingID)
		//pengelolaGedungRoutes.GET("/setting/:id", pengelolaGedungController.GetPengelolaGedungBySettingID)
		pengelolaGedungRoutes.PUT("/:id", pengelolaGedungController.UpdatePengelolaGedung)
		pengelolaGedungRoutes.DELETE("/:id", pengelolaGedungController.DeletePengelolaGedung)

	}
}
