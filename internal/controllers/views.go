package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"budget-app/internal/utils"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Views
const (
	IndexHTML        = "./views/index.html"
	SidenavHTML      = "./views/navigation/sidenav.html"
	BodyImportsHTML  = "./views/imports/bodyimports.html"
	HeadImportsHTML  = "./views/imports/headimports.html"
	DashboardHTML    = "./views/dashboard.html"
	UploadHTML       = "./views/uploads/upload.html"
	ErrorHTML        = "./views/error.html"
	Transactions     = "./views/transactions.html"
	TransactionTypes = "./views/transaction_types.html"
)

// Components
const (
	RecentTransactionComponent = "./views/components/recenttransactions.html"
	BudgetForm                 = "./views/components/budgetform.html"
)

const (
	StartDayOfMonth = "+21 days"
	DateNow         = "now"
	StringQuery     = "date(transaction_date) >= date(?, 'start of month', '-1 month', ?) AND date(transaction_date) <= date(?, 'start of month', ?)"
)

func ParseFiles(files ...string) *template.Template {
	return template.Must(template.ParseFiles(files...))
}

type PageData struct {
	Version string
}

func (ge *GinEngine) HomePage(c *gin.Context) {

	r := ge.Router

	// Generate the current timestamp in milliseconds
	version := time.Now().UnixNano() / int64(time.Millisecond)

	r.LoadHTMLFiles(IndexHTML, SidenavHTML, BodyImportsHTML, HeadImportsHTML)

	c.HTML(http.StatusOK, "homepage/index.html", gin.H{
		"Version": fmt.Sprintf("%d", version),
	})
}

func (ge *GinEngine) Dashboard(c *gin.Context) {
	r := ge.Router

	r.LoadHTMLFiles(DashboardHTML)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func ReturnUploadTemp() *template.Template {
	return ParseFiles(UploadHTML)
}

func (ge *GinEngine) UploadPage(c *gin.Context) {
	r := ge.Router

	r.LoadHTMLFiles(UploadHTML)

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func (ge *GinEngine) UploadPageRefreshed(c *gin.Context) {
	r := ge.Router

	// Generate the current timestamp in milliseconds
	version := time.Now().UnixNano() / int64(time.Millisecond)

	r.LoadHTMLFiles(IndexHTML, UploadHTML)

	c.HTML(http.StatusOK, "homepage/index.html", gin.H{
		"Version": fmt.Sprintf("%d", version),
	})
}

type BudgetItemWithTotal struct {
	Budget                 models.Budget
	TotalTransactionAmount float64
}

func (ge *GinEngine) HandleTransctions(c *gin.Context) {
	db := database.ReturnDB()
	r := ge.Router
	r.LoadHTMLFiles(RecentTransactionComponent)

	var transactions []models.Transaction
	// response := db.Joins("Budget").Find(&transactions).Scan(&transactions)
	// response := db.Where(StringQuery, DateNow, StartDayOfMonth, DateNow, StartDayOfMonth).Order("transaction_date desc").Find(&transactions).Scan(&transactions)
	response := db.Joins("Budget").Joins("TransactionType").Where(StringQuery, DateNow, StartDayOfMonth, DateNow, StartDayOfMonth).Order("transaction_date desc").Find(&transactions).Scan(&transactions)
	// response := db.Model(&transactions).Where(StringQuery, DateNow, StartDayOfMonth, DateNow, StartDayOfMonth).Order("transaction_date desc").Association("Budget").DB.Scan(&transactions)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	var budetsItems []models.Budget
	response = db.Preload("Transactions", StringQuery, DateNow, StartDayOfMonth, DateNow, StartDayOfMonth).Find(&budetsItems)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	totalIncome := 0.0
	totalExpense := 0.0
	budgetSpent := 0.0
	for _, t := range transactions {
		// if value is positive, it is income
		if t.TransactionTypeID != nil {
			if t.TransactionAmount > 0 {
				totalIncome += float64(t.TransactionAmount)
			} else {
				totalExpense += float64(t.TransactionAmount)
			}
		}
		if t.BudgetID != nil {
			budgetSpent += float64(t.TransactionAmount)
		}
	}

	recentTotal := 0.0
	for _, t := range transactions {
		recentTotal += float64(t.TransactionAmount)
	}

	budgetTotal := 0.0

	budgetTotalItems := []BudgetItemWithTotal{}

	for _, b := range budetsItems {
		budgetTotal += float64(b.Amount)
		var totalAmount float64
		for _, transaction := range b.Transactions {
			totalAmount += transaction.TransactionAmount
		}
		budgetTotalItems = append(budgetTotalItems, BudgetItemWithTotal{
			Budget:                 b,
			TotalTransactionAmount: totalAmount,
		})
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	// 	"RecentTotal":      recentTotal,
	// 	"BudgetTotal":      budgetTotal,
	// 	"TotalIncome":      totalIncome,
	// 	"TotalExpense":     totalExpense,
	// 	"Transactions":     transactions,
	// 	"TransactionCount": len(transactions),
	// 	"BudgetItems":      budetsItems,
	// 	"BudgetSpent":      budgetSpent,
	// 	"BudgetTotalItems": budgetTotalItems,
	// })

	c.HTML(http.StatusOK, "recenttransactions.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"RecentTotal":      recentTotal,
		"BudgetTotal":      budgetTotal,
		"TotalIncome":      totalIncome,
		"TotalExpense":     totalExpense,
		"Transactions":     transactions,
		"TransactionCount": len(transactions),
		"BudgetItems":      budetsItems,
		"BudgetSpent":      budgetSpent,
		"BudgetTotalItems": budgetTotalItems,
	})
}

type TransactionData struct {
	Transaction      models.Transaction
	TransactionTypes []models.TransactionType
	BudgetItems      []models.Budget
}

func (ge *GinEngine) ReturnTransactions(c *gin.Context) {
	r := ge.Router
	db := database.ReturnDB()

	r.LoadHTMLFiles(Transactions)

	var transactions []models.Transaction
	response := db.Preload("Budget").Order("transaction_date desc").Find(&transactions).Scan(&transactions)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
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

	var budgetItems []models.Budget
	response = db.Find(&budgetItems).Scan(&budgetItems)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	data := []TransactionData{}

	for _, transaction := range transactions {
		data = append(data, TransactionData{
			Transaction:      transaction,
			TransactionTypes: transactionTypes,
			BudgetItems:      budgetItems,
		})
	}

	c.HTML(http.StatusOK, "transactions.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"TransactionsData": data,
		"TransactionCount": len(transactions),
	})
}

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

func (ge *GinEngine) TransactionsAddTransactionType(c *gin.Context) {
	r := ge.Router
	r.LoadHTMLFiles(Transactions)
	db := database.ReturnDB()

	transactionID := c.Param("id")
	transactionTypeID := c.PostForm("transactionTypeID")

	var transaction models.Transaction
	transaction.ID = utils.StringToUint(transactionID)

	response := db.Model(&transaction).Update("transaction_type_id", transactionTypeID).Scan(&transaction)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not add budget item to transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"transaction": transaction,
	})
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
