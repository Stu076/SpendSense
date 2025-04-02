package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Stu076/SpendSense/internal/cache"
	"github.com/Stu076/SpendSense/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type ExpenseHandler struct {
	DB *bun.DB
}

// TODO: add UPDATE expense and get expense by ID

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	username, exists := c.Get("username")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var expense models.Expense
	if err := c.ShouldBindJSON(&expense); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User
	err := h.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	expense.UserID = user.ID

	if err := models.CreateExpense(h.DB, &expense); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create expense"})
		return
	}

	// Delete from Cache
	err = cache.DeleteCache("expenses:" + strconv.Itoa(user.ID))
	if err != nil {
		log.Println("❌ Failed to delete cache:", err)
	}
	// TODO: set cache for the new expense

	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "expense": expense})
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	username, exists := c.Get("username")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var user models.User
	err := h.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	cacheKey := "expenses:" + strconv.Itoa(user.ID)
	cachedExpenses, err := cache.GetCache(cacheKey)

	// Cache hit
	if err == nil {
		log.Println("✅ Serving from Redis Cache")
		var expenses []models.Expense
		err = json.Unmarshal([]byte(cachedExpenses), &expenses)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse cached expenses"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"expenses": expenses})
		return
	}

	// Cache miss
	var expenses []models.Expense
	err = h.DB.NewSelect().Model(&expenses).Where("user_id = ?", user.ID).Order("created_at DESC").Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	// Set cache
	expensesJSON, err := json.Marshal(expenses)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache expenses"})
		return
	}

	err = cache.SetCache(cacheKey, string(expensesJSON), 10*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cache expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	username, exists := c.Get("username")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	expenseID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expense ID"})
		return
	}

	// Get User ID
	var user models.User
	err = h.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Delete from DB
	_, err = h.DB.NewDelete().Model((*models.Expense)(nil)).Where("id = ? AND user_id = ?", expenseID, user.ID).Exec(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete expense"})
		return
	}

	// Delete from Cache
	err = cache.DeleteCache("expenses:" + strconv.Itoa(user.ID))
	if err != nil {
		log.Println("❌ Failed to delete cache:", err)
	}
	// TODO: delete cache for the deleted expense

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
}
