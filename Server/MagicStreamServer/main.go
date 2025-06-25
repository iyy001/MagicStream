package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"strings"

	"github.com/GavinLonDigital/MagicStream/Server/MagicStreamServer/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Set Gin to release mode to avoid debug logs in production
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: .env file not found")
	}
}

func main() {
	// Create a new Gin router
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, MagicStream!")
	})

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
		log.Fatal("PORT not set in .env file")
	}
	fmt.Println("MongoDB URI:", port)

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	routes.SetupUnprotectedRoutes(router)
	routes.SetupProtectedRoutes(router)
	// Start the server on the specified port
	if err := router.Run(":" + port); err != nil {
		fmt.Println("Failed to start server:", err)
	}

}
