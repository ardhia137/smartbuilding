package infrastructure

import (
	"github.com/gin-gonic/gin"
	"smartbuilding/controllers"
)

func RegisterUserRoutes(router *gin.Engine, userController *controllers.UserController) {
	apiGroup := router.Group("/api") // Tambahkan prefix /api
	userRoutes := apiGroup.Group("/users")
	{
		userRoutes.GET("", userController.GetAllUsers)
		userRoutes.POST("", userController.CreateUser)
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}
}
func RegisterKamarRoutes(router *gin.Engine, kamarController *controllers.KamarController) {
	apiGroup := router.Group("/api") // Tambahkan prefix /api

	kamarRoutes := apiGroup.Group("/kamar")
	{
		kamarRoutes.GET("", kamarController.GetAllKamars)
		kamarRoutes.POST("", kamarController.CreateKamar)
		kamarRoutes.GET("/:id", kamarController.GetKamarByID)
		kamarRoutes.PUT("/:id", kamarController.UpdateKamar)
		kamarRoutes.DELETE("/:id", kamarController.DeleteKamar)
	}
}
