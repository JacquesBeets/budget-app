package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBudget(s database.Service) ([]models.Budget, error) {
	db := s.GetDBPool()

	rows, err := db.Query(`
		SELECT id, name, amount, created_at, transaction_type_id
		FROM budget;
	`)

	if err != nil {
		fmt.Print("Error getting budget:", err)
		return nil, err
	}
	defer rows.Close()

	var budget []models.Budget
	for rows.Next() {
		var b models.Budget
		if err := rows.Scan(&b.ID, &b.Name, &b.Amount, &b.CreatedAt, &b.TransactionTypeID); err != nil {
			log.Fatalf(fmt.Sprintf("Error scanning budget: %v", err))
			return nil, err
		}
		budget = append(budget, b)
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
