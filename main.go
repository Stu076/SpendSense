package main

import (
	"log"

	"github.com/Stu076/SpendSense/config"
	"github.com/Stu076/SpendSense/internal/cache"
	"github.com/Stu076/SpendSense/internal/handlers"
	"github.com/Stu076/SpendSense/internal/middlewares"
	"github.com/Stu076/SpendSense/internal/repositories"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

func main() {
	confStatus := config.Init()
	if !confStatus {
		panic("Error initializing config")
	}

	cache.InitRedis()

	db := repositories.MustSetupDatabase()

	r := gin.Default()

	// Add Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Authentication
	authHandler := handlers.AuthHandler{DB: db}
	r.POST("/register", middlewares.TrackRequest("POST", "/register", authHandler.RegisterUser))
	r.POST("/login", middlewares.TrackRequest("POST", "/login", authHandler.LoginUser))

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())

	// Expenses
	expenseHandler := handlers.ExpenseHandler{DB: db}
	protected.POST("/expenses", middlewares.TrackRequest("POST", "/api/expenses", expenseHandler.CreateExpense))
	protected.GET("/expenses", middlewares.TrackRequest("GET", "/api/expenses", expenseHandler.GetExpenses))
	protected.DELETE("/expenses/:id", middlewares.TrackRequest("DELETE", "/api/expenses/:id", expenseHandler.DeleteExpense))

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
