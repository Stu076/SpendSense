package main

import (
	"log"
	"net/http"

	"github.com/Stu076/SpendSense/config"
	"github.com/Stu076/SpendSense/internal/handlers"
	"github.com/Stu076/SpendSense/internal/middlewares"
	"github.com/Stu076/SpendSense/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	confStatus := config.Init()
	if !confStatus {
		panic("Error initializing config")
	}

	db := repositories.MustSetupDatabase()

	r := gin.Default()

	authHandler := handlers.AuthHandler{DB: db}
	r.POST("/register", authHandler.RegisterUser)
	r.POST("/login", authHandler.LoginUser)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/secure", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "You are authenticated!"})
	})
	// Load Config
	// viper.SetConfigFile(".env")
	// viper.AutomaticEnv()
	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Println("No config file found, using environment variables")
	// }

	// // Initialize Gin Router
	// r := gin.Default()

	// // Test route
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "pong"})
	// })

	// Start server
	port := viper.GetString("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Starting server on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
