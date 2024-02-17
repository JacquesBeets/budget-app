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
	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
		"isEmpty":     utils.IsEmpty,
		"isNil":       utils.IsNil,
		"isTotalSpendGreaterThanBudget": utils.IsTotalSpendGreaterThanBudget,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(RecentTransactionComponent)

	var transactions []models.Transaction
	response := db.Where("date(transaction_date) >= date(?, 'start of month', '-1 month', '+21 days')", "now").Order("transaction_date desc").Find(&transactions).Scan(&transactions)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	var budetsItems []models.Budget
	response = db.Preload("Transactions").Find(&budetsItems)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	// Print budget items along with their transactions
	// for _, b := range budetsItems {
	// 	fmt.Printf("Budget: %s (Amount: %f)\n", b.Name, b.Amount)
	// 	for _, t := range b.Transactions {
	// 		fmt.Printf("Transaction: Type: %s, Amount: %f\n", t.BankTransactionType, t.TransactionAmount)
	// 	}
	// }

	totalIncome := 0.0
	totalExpense := 0.0
	budgetSpent := 0.0
	for _, t := range transactions {
		// if value is positive, it is income
		if t.TransactionAmount > 0 {
			totalIncome += float64(t.TransactionAmount)
		} else {
			totalExpense += float64(t.TransactionAmount)
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

	funcMap := template.FuncMap{
		"formatDate":      utils.FormatDate,
		"formatPrice":     utils.FormatPrice,
		"isEmpty":         utils.IsEmpty,
		"isNil":           utils.IsNil,
		"typeOf":          utils.CheckType,
		"dereferencePntr": utils.DereferenceUintPtr,
	}
	r.SetFuncMap(funcMap)
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
	// r := ge.Router

	// funcMap := template.FuncMap{
	// 	"formatDate":  utils.FormatDate,
	// 	"formatPrice": utils.FormatPrice,
	// }
	// r.SetFuncMap(funcMap)
	// r.LoadHTMLFiles(TransactionTypes)

	// service := database.New()
	// transactionTypes, err := GetTransactionsTypes(service)
	// if err != nil {
	// 	r.LoadHTMLFiles(ErrorHTML)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction types"})
	// 	return
	// }

	// c.HTML(http.StatusOK, "transaction_types.html", gin.H{
	// 	"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	// 	"TransactionTypes": transactionTypes,
	// 	"TransactionCount": len(transactionTypes),
	// })
}

func (ge *GinEngine) HandleTransactionTypeCreate(c *gin.Context) {
	// r := ge.Router
	// funcMap := template.FuncMap{
	// 	"formatDate":  utils.FormatDate,
	// 	"formatPrice": utils.FormatPrice,
	// }
	// r.SetFuncMap(funcMap)
	// r.LoadHTMLFiles(TransactionTypes)

	// service := database.New()
	// transactionType := models.TransactionType{}
	// transactionType.Title = c.PostForm("title")
	// transactionType.Category = c.PostForm("category")

	// transactionType, err := CreateTransactionType(service, transactionType)
	// if err != nil {
	// 	r.LoadHTMLFiles(ErrorHTML)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create transaction type"})
	// 	return
	// }

	// transactionTypes, err := GetTransactionsTypes(service)
	// if err != nil {
	// 	r.LoadHTMLFiles(ErrorHTML)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch transaction types"})
	// 	return
	// }

	// c.HTML(http.StatusOK, "transaction_types.html", gin.H{
	// 	"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	// 	"TransactionTypes": transactionTypes,
	// 	"TransactionCount": 1,
	// })
}

func (ge *GinEngine) TransactionsAddTransactionType(c *gin.Context) {
	// r := ge.Router
	// funcMap := template.FuncMap{
	// 	"formatDate":  utils.FormatDate,
	// 	"formatPrice": utils.FormatPrice,
	// }
	// r.SetFuncMap(funcMap)
	// r.LoadHTMLFiles(Transactions)

	// service := database.New()
	// transactionID := c.Param("id")
	// transactionTypeID := c.PostForm("transactionTypeID")

	// err := LinkTransactionType(service, transactionID, transactionTypeID)

	// if err != nil {
	// 	r.LoadHTMLFiles(ErrorHTML)
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "could not add transaction type to transaction",
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"status": "ok",
	// })
}

func (ge *GinEngine) BudgetTransactionAdd(c *gin.Context) {
	r := ge.Router
	db := database.ReturnDB()
	// funcMap := template.FuncMap{
	// 	"formatDate":  utils.FormatDate,
	// 	"formatPrice": utils.FormatPrice,
	// }
	// r.SetFuncMap(funcMap)
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
