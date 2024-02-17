package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (ge *GinEngine) ReturnTransactionTypes(c *gin.Context) {
	r := ge.Router
	r.LoadHTMLFiles(TransactionTypes)

	db := database.ReturnDB()
	var transactionTypes []models.TransactionType
	response := db.Find(&transactionTypes).Scan(&transactionTypes)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
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

	db := database.ReturnDB()

	category := c.PostForm("category")
	transactionType := &models.TransactionType{
		Title:    c.PostForm("title"),
		Category: &category,
	}

	response := db.Create(transactionType)
	if response.Error != nil {
		fmt.Println("Error saving budget: ", response.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save budget"})
		return
	}

	var transactionTypes []models.TransactionType
	response = db.Find(&transactionTypes).Scan(&transactionTypes)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	c.HTML(http.StatusOK, "transaction_types.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"TransactionTypes": transactionTypes,
		"TransactionCount": 1,
	})
}
