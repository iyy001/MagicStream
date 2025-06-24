package routes

import (
	controller "github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnprotectedRoutes(router *gin.Engine) {
	router.POST("/register", controller.RegisterUser())
	router.POST("/login", controller.LoginUser())
	router.GET("/movies", controller.GetMovies())
	router.GET("/genres", controller.GetGenres())
	router.POST("/logout", controller.LogoutHandler())
	router.POST("/refresh", controller.RefreshTokenHandler())
}
