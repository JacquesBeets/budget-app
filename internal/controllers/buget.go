package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"budget-app/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ge *GinEngine) SaveBudgetItem(c *gin.Context) {
	dbService := database.ReturnDB()
	var budget *models.Budget

	name := c.PostForm("name")
	amountStr := c.PostForm("amount")

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error converting amount to float64: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	budget = &models.Budget{
		Name:   name,
		Amount: amount,
	}

	response := dbService.Create(budget).Scan(&budget)
	if response.Error != nil {
		fmt.Println("Error saving budget: ", response.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (ge *GinEngine) BudgetTransactionAdd(c *gin.Context) {
	r := ge.Router
	db := database.ReturnDB()
	// r.LoadHTMLFiles(Transactions)

	transactionID := c.Param("id")
	budgetID := c.PostForm("budgetItemID")

	var transaction models.Transaction
	transaction.ID = utils.StringToUint(transactionID)

	response := db.Model(&transaction).Update("budget_id", budgetID).Scan(&transaction)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not add budget item to transaction",
		})
		return
	}

	fmt.Println("Transaction: ", transaction)

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"transaction": transaction,
	})
}

func (ge *GinEngine) ReturnBudgetForm(c *gin.Context) {
	// single file
	r := ge.Router
	r.LoadHTMLFiles(BudgetForm)

	c.HTML(http.StatusOK, "budgetform.html", gin.H{
		"Version": "1",
	})
}
