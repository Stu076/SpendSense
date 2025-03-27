package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Stu076/SpendSense/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type ExpenseHandler struct {
	DB *bun.DB
}

func (h *ExpenseHandler) CreateExpense(c *gin.Context) {
	username, _ := c.Get("username")

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

	c.JSON(http.StatusCreated, gin.H{"message": "Expense created successfully", "expense": expense})
}

func (h *ExpenseHandler) GetExpenses(c *gin.Context) {
	username, _ := c.Get("username")

	var user models.User
	err := h.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var expenses []models.Expense
	err = h.DB.NewSelect().Model(&expenses).Where("user_id = ?", user.ID).Order("created_at DESC").Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch expenses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"expenses": expenses})
}

func (h *ExpenseHandler) DeleteExpense(c *gin.Context) {
	username, _ := c.Get("username")
	expenseID, _ := strconv.Atoi(c.Param("id"))

	// Get User ID
	var user models.User
	err := h.DB.NewSelect().Model(&user).Where("username = ?", username).Scan(context.Background())
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

	c.JSON(http.StatusOK, gin.H{"message": "Expense deleted"})
}
