package controllers

import (
	"budget-app/internal/models"
	"budget-app/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (ge *GinEngine) SaveBudgetItem(c *gin.Context) {
	var budget *models.Budget

	name := c.PostForm("name")
	amountStr := c.PostForm("amount")

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	budget = &models.Budget{
		Name:   name,
		Amount: amount,
	}

	response := ge.db().Create(budget).Scan(&budget)
	if response.Error != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	c.JSON(http.StatusOK, budget)
}

func (ge *GinEngine) BudgetTransactionAdd(c *gin.Context) {

	transactionID := c.Param("id")
	budgetID := utils.StringToUint(c.PostForm("budgetItemID"))

	var transaction models.Transaction
	transaction.ID = utils.StringToUint(transactionID)

	response := ge.db().Model(&transaction).Update("budget_id", &budgetID).Scan(&transaction)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
		return
	}

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
