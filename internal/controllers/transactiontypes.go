package controllers

import (
	"budget-app/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (ge *GinEngine) ReturnTransactionTypes(c *gin.Context) {
	r := ge.Router
	r.LoadHTMLFiles(TransactionTypes)
	var transactionTypes []models.TransactionType
	response := ge.db().Find(&transactionTypes).Scan(&transactionTypes)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
		return
	}

	c.HTML(http.StatusOK, "transaction_types.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"TransactionTypes": transactionTypes,
		"TransactionCount": len(transactionTypes),
	})
}

func (ge *GinEngine) HandleTransactionTypeCreate(c *gin.Context) {
	r := ge.Router
	r.LoadHTMLFiles(TransactionTypes)

	category := c.PostForm("category")
	transactionType := &models.TransactionType{
		Title:    c.PostForm("title"),
		Category: &category,
	}

	response := ge.db().Create(transactionType)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
		return
	}

	var transactionTypes []models.TransactionType
	response = ge.db().Find(&transactionTypes).Scan(&transactionTypes)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
		return
	}

	c.HTML(http.StatusOK, "transaction_types.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"TransactionTypes": transactionTypes,
		"TransactionCount": 1,
	})
}
