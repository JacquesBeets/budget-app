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

	// tmpl := ParseFiles(DashboardHTML)

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

func (ge *GinEngine) HandleTransctions(c *gin.Context) {
	service := database.New()
	r := ge.Router
	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(RecentTransactionComponent)

	transactions, err := GetTransactions(service)
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "views/error.html", gin.H{"error": "could not fetch transactions"})
		return
	}

	budetsItems, err := GetBudget(service)
	if err != nil {
		fmt.Println("Error getting budget items: ", err)
	}

	totalIncome := 0.0
	totalExpense := 0.0
	for _, t := range transactions {
		// if value is positive, it is income
		if t.TransactionAmount > 0 {
			totalIncome += float64(t.TransactionAmount)
		} else {
			totalExpense += float64(t.TransactionAmount)
		}
	}

	recentTotal := 0.0
	for _, t := range transactions {
		recentTotal += float64(t.TransactionAmount)
	}

	budgetTotal := 0.0
	for _, b := range budetsItems {
		budgetTotal += float64(b.Amount)
	}

	c.HTML(http.StatusOK, "recenttransactions.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"RecentTotal":      recentTotal,
		"BudgetTotal":      budgetTotal,
		"TotalIncome":      totalIncome,
		"TotalExpense":     totalExpense,
		"Transactions":     transactions,
		"TransactionCount": len(transactions),
		"BudgetItems":      budetsItems,
	})
}

func (ge *GinEngine) ReturnTransactions(c *gin.Context) {
	r := ge.Router
	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(Transactions)

	service := database.New()
	transactions, err := GetAllTransactions(service)
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transactions"})
		return
	}

	transactionTypes, err := GetTransactionsTypes(service)
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction types"})
		return
	}

	type TransactionData struct {
		Transaction      models.Transaction
		TransactionTypes []models.TransactionType
	}

	data := []TransactionData{}

	for _, transaction := range transactions {
		data = append(data, TransactionData{
			Transaction:      transaction,
			TransactionTypes: transactionTypes,
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

	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(TransactionTypes)

	service := database.New()
	transactionTypes, err := GetTransactionsTypes(service)
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction types"})
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
	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(TransactionTypes)

	service := database.New()
	transactionType := models.TransactionType{}
	transactionType.Title = c.PostForm("title")
	transactionType.Category = c.PostForm("category")

	transactionType, err := CreateTransactionType(service, transactionType)
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create transaction type"})
		return
	}

	transactionTypes, err := GetTransactionsTypes(service)
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction types"})
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
	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(Transactions)

	service := database.New()
	transactionID := c.Param("id")
	transactionTypeID := c.PostForm("transactionTypeID")

	err := LinkTransactionType(service, transactionID, transactionTypeID)

	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not add transaction type to transaction",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
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
