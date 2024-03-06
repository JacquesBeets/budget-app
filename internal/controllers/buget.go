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

func (ge *GinEngine) UpdateBudgetItem(c *gin.Context) {
	var budget *models.Budget

	budgetID := c.Param("id")
	name := c.PostForm("name")
	amountStr := c.PostForm("amount")

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	// Find the budget item
	result := ge.db().First(&budget, utils.StringToUint(budgetID))
	if result.Error != nil {
		ge.ReturnErrorJSON(c, result.Error)
		return
	}

	// Update the budget item
	budget.Name = name
	budget.Amount = amount
	response := ge.db().Model(&budget).Updates(budget).Scan(&budget)
	if response.Error != nil {
		ge.ReturnErrorJSON(c, response.Error)
		return
	}

	ge.HandleTransctions(c)
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

func (ge *GinEngine) ReturnBudgetEditForm(c *gin.Context) {
	// single file
	r := ge.Router
	r.LoadHTMLFiles(BudgetEdit)
	budgetID := c.Param("id")

	var budget models.Budget
	response := ge.db().First(&budget, budgetID)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
		return
	}

	c.HTML(http.StatusOK, "budgetedit.html", gin.H{
		"Budget": budget,
	})
}
