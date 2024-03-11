package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"budget-app/internal/utils"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aclindsa/ofxgo"
	"github.com/gin-gonic/gin"
)

func (ge *GinEngine) HandleOFXUpload(c *gin.Context) {
	// single file
	r := ge.Router
	file, _ := c.FormFile("ofx")
	bank := c.PostForm("bank")
	dst := "./uploads/" + bank + "-" + file.Filename

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dst)

	r.LoadHTMLFiles(UploadHTML)

	err := ParseOFX(dst, bank)
	if err != nil {
		fmt.Printf("could not parse OFX file: %v", err)
		return
	}

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func ParseOFX(filePath, bankName string) error {
	db := database.ReturnDB()

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("could not open OFX file: %v", err)
		return err
	}
	defer f.Close()

	resp, err := ofxgo.ParseResponse(f)
	if err != nil {
		fmt.Printf("could not parse OFX file: %v", err)
		return err
	}

	var Transactions models.Transactions

	// Access the Bank Messages
	if len(resp.Bank) > 0 {
		bankMessage := resp.Bank[0]
		if stmt, ok := bankMessage.(*ofxgo.StatementResponse); ok {
			// Access the TransactionList
			transactions := stmt.BankTranList
			for _, transaction := range transactions.Transactions {

				// Create Transaction
				amount, _ := transaction.TrnAmt.Float64()
				trnAmt := float64(amount)

				var transactionModel models.Transaction

				trn, err := transactionModel.New(
					fmt.Sprint(transaction.TrnType),
					transaction.DtPosted.Format("2006-01-02"),
					trnAmt,
					string(transaction.FiTID),
					string(transaction.Name),
					string(transaction.Memo),
					bankName,
				)
				if err != nil {
					fmt.Printf("could not create transaction: %v", err)
					return err
				}

				// Check if transaction already exists
				response := trn.Exists(db)
				if !response {
					trn.AutoCategorize(db)
					Transactions = append(Transactions, *trn)
				} else {
					log.Println("Transaction already exists")
					continue
				}

			}

			if len(Transactions) > 0 {
				createdTransactions := Transactions.Create(db)
				// Transactions.PrintAll()
				if createdTransactions.Error != nil {
					fmt.Printf("could not create transactions: %v", createdTransactions.Error)
					return createdTransactions.Error
				}
			}
		}
	}

	return nil
}

type BudgetItemWithTotal struct {
	Budget                 models.Budget
	TotalTransactionAmount float64
}

func (ge *GinEngine) HandleTransctions(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		date = "now"
	} else {
		date = date + "-01"
	}

	r := ge.Router
	r.LoadHTMLFiles(RecentTransactionComponent)

	var transactions []models.Transaction
	response := ge.db().Joins("Budget").Joins("TransactionType").Where(StringQuery, date, StartDayOfMonth, date, StartDayOfMonth).Order("transaction_date desc").Find(&transactions).Scan(&transactions)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	var budetsItems []models.Budget
	response = ge.db().Preload("Transactions", StringQuery, date, StartDayOfMonth, date, StartDayOfMonth).Order("amount desc").Find(&budetsItems)
	if response.Error != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not fetch response"})
		fmt.Println("Error getting response: ", response.Error)
		return
	}

	totalIncome := 0.0
	totalExpense := 0.0
	budgetSpent := 0.0
	recentTotal := 0.0
	for _, t := range transactions {
		if t.TransactionTypeID != nil && *t.TransactionTypeID != 12 {
			if t.TransactionAmount > 0 {
				// fmt.Println(*t.TransactionTypeID, "Adding to totalIncome:", t.TransactionAmount)
				totalIncome += float64(t.TransactionAmount)
			} else {
				// fmt.Println(*t.TransactionTypeID, "Adding to totalExpense:", t.TransactionAmount)
				totalExpense += float64(t.TransactionAmount) // Adjusting for negative expenses
			}
		}
		if t.BudgetID != nil {
			// fmt.Println("Adding to budgetSpent:", t.TransactionAmount)
			budgetSpent += float64(t.TransactionAmount)
		}
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
		"Balance":          totalIncome + totalExpense,
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
	date := c.Query("date") // 2024-02
	if date == "" {
		date = "now"
	} else {
		date = date + "-01"
	}

	r := ge.Router
	r.LoadHTMLFiles(Transactions)

	var transactions []models.Transaction
	response := ge.db().Preload("Budget").Where(StringQuery, date, StartDayOfMonth, date, StartDayOfMonth).Order("transaction_date desc").Find(&transactions).Scan(&transactions)
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

	var budgetItems []models.Budget
	response = ge.db().Find(&budgetItems).Scan(&budgetItems)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
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

func (ge *GinEngine) TransactionsAddTransactionType(c *gin.Context) {
	r := ge.Router
	r.LoadHTMLFiles(Transactions)

	var transaction models.Transaction

	transactionIDUint := utils.StringToUint(c.Param("id"))
	transactionTypeIDUint := utils.StringToUint(c.PostForm("transactionTypeID"))

	transaction.ID = transactionIDUint
	transaction.TransactionTypeID = &transactionTypeIDUint

	response := ge.db().Model(&transaction).Update("transaction_type_id", &transactionTypeIDUint).Scan(&transaction)
	if response.Error != nil {
		ge.ReturnErrorPage(c, response.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ok",
		"transaction": transaction,
	})
}
