package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBudget() ([]models.Budget, error) {
	var dbService database.Service = database.New()
	budget, err := dbService.GetBudget()
	if err != nil {
		return nil, err
	}
	return budget, nil
}

func NewBudget(
	name string,
	amount float64,
	transactionTypeID int,
) (*models.Budget, error) {

	return &models.Budget{
		Name:              name,
		Amount:            amount,
		TransactionTypeID: transactionTypeID,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func (ge *GinEngine) SaveBudgetItem(c *gin.Context) {
	var dbService database.Service = database.New()
	var budget *models.Budget

	name := c.PostForm("name")
	amountStr := c.PostForm("amount")
	transactionTypeIDStr := "1"

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	// Convert transactionTypeID to int
	transactionTypeID, err := strconv.Atoi(transactionTypeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction type ID"})
		return
	}

	budget, err = NewBudget(name, amount, transactionTypeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating budget"})
		return
	}

	err = dbService.SaveBudgetItem(*budget)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error saving budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}
