package infrastructure

import (
	"smartbuilding/controllers"
	"smartbuilding/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	apiGroup := router.Group("/api")

	adminRoutes := apiGroup.Group("/users")
	adminRoutes.Use(utils.RoleMiddleware("admin", "manajement"), utils.UserIDMiddleware())
	{
		adminRoutes.GET("", userController.GetAllUsers)
		adminRoutes.POST("", userController.CreateUser)
		adminRoutes.GET("/:id", userController.GetUserByID)
		adminRoutes.PUT("/:id", userController.UpdateUser)
		adminRoutes.DELETE("/:id", userController.DeleteUser)
	}

	userSelfRoutes := apiGroup.Group("/users")
	userSelfRoutes.Use(utils.UserIDMiddleware())
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

func RegisterGedungRoutes(router *gin.Engine, gedungController *controllers.GedungController) {
	apiGroup := router.Group("/api")
	gedungRoutes := apiGroup.Group("/gedung")
	gedungRoutes.Use(utils.RoleMiddleware("admin", "manajement", "pengelola"), utils.UserIDMiddleware())
	{
		gedungRoutes.GET("", gedungController.GetAllGedung)
		gedungRoutes.POST("", gedungController.CreateGedung)
		gedungRoutes.GET("/:id", gedungController.GetGedungByID)
		gedungRoutes.PUT("/:id", gedungController.UpdateGedung)
		gedungRoutes.DELETE("/:id", gedungController.DeleteGedung)
	}
}

func RegisterTorentRoutes(router *gin.Engine, torentController *controllers.TorentController) {
	apiGroup := router.Group("/api")
	torentRoutes := apiGroup.Group("/torent")
	torentRoutes.Use(utils.RoleMiddleware("admin", "manajement", "pengelola"), utils.UserIDMiddleware())
	{
		torentRoutes.GET("", torentController.GetAllTorent)
		torentRoutes.POST("", torentController.CreateTorent)
		torentRoutes.GET("/:id", torentController.GetTorentByID)
		torentRoutes.GET("/gedung/:id", torentController.GetTorentByGedungID)
		torentRoutes.PUT("/:id", torentController.UpdateTorent)
		torentRoutes.DELETE("/:id", torentController.DeleteTorent)
	}
}

func RegisterHakAksesRoutes(router *gin.Engine, hakAksesController *controllers.HakAksesController) {
	apiGroup := router.Group("/api")
	hakAksesRoutes := apiGroup.Group("/hak_akses")
	hakAksesRoutes.Use(utils.RoleMiddleware("admin", "manajement"), utils.UserIDMiddleware())
	{
		hakAksesRoutes.GET("", hakAksesController.GetAllHakAkses)
		hakAksesRoutes.POST("", hakAksesController.CreateHakAkses)
		hakAksesRoutes.GET("/:id", hakAksesController.GetHakAksesByID)
		hakAksesRoutes.GET("/gedung/:id", hakAksesController.GetHakAksesByGedungID)
		hakAksesRoutes.PUT("/:id", hakAksesController.UpdateHakAkses)
		hakAksesRoutes.DELETE("/:id", hakAksesController.DeleteHakAkses)

	}
}
