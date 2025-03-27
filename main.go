package main

import (
	"log"

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

	// Authentication
	authHandler := handlers.AuthHandler{DB: db}
	r.POST("/register", authHandler.RegisterUser)
	r.POST("/login", authHandler.LoginUser)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthMiddleware())

	// Expenses
	expenseHandler := handlers.ExpenseHandler{DB: db}
	protected.POST("/expenses", expenseHandler.CreateExpense)
	protected.GET("/expenses", expenseHandler.GetExpenses)
	protected.DELETE("/expenses/:id", expenseHandler.DeleteExpense)

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
