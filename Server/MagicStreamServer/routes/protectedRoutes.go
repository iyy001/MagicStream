package routes

import (
	controller "github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/controllers"
	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	router.Use(middleware.AuthMiddleware())
	router.GET("/movie/:imdb_id", controller.GetMovie())
	router.GET("/recommendedmovies", controller.GetRecommendedMovies())
	router.POST("/addmovie", controller.AddMovie())
	router.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate())
}
