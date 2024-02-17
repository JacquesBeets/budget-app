package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aclindsa/ofxgo"
	"github.com/gin-gonic/gin"
)

func NewTransaction(
	transactionType string,
	transactionDate string,
	transactionAmount float64,
	transactionID string,
	transactionName string,
	transactionMemo string,
	bankName string,
) (*models.Transaction, error) {

	transactionDateParsed, err := time.Parse("2006-01-02", transactionDate)
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		BankTransactionType: &transactionType,
		TransactionDate:     transactionDateParsed,
		TransactionAmount:   transactionAmount,
		BankTransactionID:   transactionID,
		TransactionName:     &transactionName,
		TransactionMemo:     &transactionMemo,
		BankName:            bankName,
		CreatedAt:           time.Now().UTC(),
	}, nil
}

// func GetAllTransactions(s database.Service) ([]models.Transaction, error) {
// 	transactions := []models.Transaction{}

// 	query := `SELECT * FROM transactions ORDER BY date(transaction_date) DESC;`

// 	rows, err := s.Query(query)
// 	if err != nil {
// 		fmt.Println(err)
// 		return transactions, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var transaction models.Transaction
// 		var transactionDateString string

// 		err := rows.Scan(
// 			&transaction.ID,
// 			&transaction.TransactionType,
// 			&transactionDateString,
// 			&transaction.TransactionAmount,
// 			&transaction.TransactionID,
// 			&transaction.TransactionName,
// 			&transaction.TransactionMemo,
// 			&transaction.CreatedAt,
// 			&transaction.BankName,
// 			&transaction.TransactionTypeID,
// 			&transaction.BudgetItemID,
// 		)
// 		if err != nil {
// 			fmt.Println(err)
// 			return transactions, err
// 		}

// 		// Parse the date string into a time.Time type
// 		transaction.TransactionDate, err = time.Parse("2006-01-02 15:04:05-07:00", transactionDateString)
// 		if err != nil {
// 			fmt.Println(err)
// 			return transactions, err
// 		}
// 		transactions = append(transactions, transaction)
// 	}

// 	return transactions, nil
// }

// func GetTransactions(s database.Service) ([]models.Transaction, error) {
// 	transactions := []models.Transaction{}
// 	currentMonth := time.Now().Format("01") // "01" is the format for two-digit month in Go

	// query := `SELECT *
	// FROM transactions
	// WHERE date(transaction_date) >= date('now', 'start of month', '-1 month', '+23 days')
	// ORDER BY date(transaction_date) DESC;`

// 	rows, err := s.Query(query, currentMonth)
// 	if err != nil {
// 		fmt.Println(err)
// 		return transactions, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var transaction models.Transaction
// 		var transactionDateString string

// 		err := rows.Scan(
// 			&transaction.ID,
// 			&transaction.TransactionType,
// 			&transactionDateString,
// 			&transaction.TransactionAmount,
// 			&transaction.TransactionID,
// 			&transaction.TransactionName,
// 			&transaction.TransactionMemo,
// 			&transaction.CreatedAt,
// 			&transaction.BankName,
// 			&transaction.TransactionTypeID,
// 			&transaction.BudgetItemID,
// 		)
// 		if err != nil {
// 			fmt.Println(err)
// 			return transactions, err
// 		}

// 		// Parse the date string into a time.Time type
// 		transaction.TransactionDate, err = time.Parse("2006-01-02 15:04:05-07:00", transactionDateString)
// 		if err != nil {
// 			fmt.Println(err)
// 			return transactions, err
// 		}
// 		transactions = append(transactions, transaction)
// 	}

// 	return transactions, nil
// }

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
		// c.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": "could not parse OFX file"})
		return
	}

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func ParseOFX(filePath, bankName string) error {
	dbService := database.ReturnDB()

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

	var Transactions []models.Transaction

	// Access the Bank Messages
	if len(resp.Bank) > 0 {
		bankMessage := resp.Bank[0]
		if stmt, ok := bankMessage.(*ofxgo.StatementResponse); ok {
			// Access the TransactionList
			transactions := stmt.BankTranList
			for _, transaction := range transactions.Transactions {

				// Add Transaction
				amount, _ := transaction.TrnAmt.Float64()
				trnAmt := float64(amount)

				trn, err := NewTransaction(
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
				Transactions = append(Transactions, *trn)
			}

			createdTransactions := dbService.Create(Transactions)

			if createdTransactions.Error != nil {
				fmt.Printf("could not create transactions: %v", createdTransactions.Error)
				return createdTransactions.Error
			}
		}
	}

	return nil
}

// func SaveMultipleTransactions(s database.Service, transactions []models.Transaction) error {
// 	db := s.GetDBPool()

// 	// Prepare the query for adding transactions
// 	query := `insert into transactions (
//         transaction_type,
//         transaction_date,
//         transaction_amount,
//         transaction_id,
//         transaction_name,
//         transaction_memo,
//         created_at,
//         transaction_type_id,
//         bank_name
//     ) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`

// 	stmt, err := db.Prepare(query)
// 	if err != nil {
// 		return err
// 	}
// 	defer stmt.Close()

// 	for _, transaction := range transactions {

// 		// Check if the transaction already exists
// 		existsQuery := `select count(*) from transactions where transaction_id = ?`
// 		var count int
// 		err = db.QueryRow(existsQuery, transaction.TransactionID).Scan(&count)
// 		if err != nil {
// 			return err
// 		}
// 		if count > 0 {
// 			// Transaction already exists, do not add it
// 			log.Println("Transaction already exists")
// 			continue
// 		}

// 		_, err = stmt.Exec(
// 			transaction.TransactionType,
// 			transaction.TransactionDate,
// 			transaction.TransactionAmount,
// 			transaction.TransactionID,
// 			transaction.TransactionName,
// 			transaction.TransactionMemo,
// 			transaction.CreatedAt,
// 			transaction.TransactionTypeID,
// 			transaction.BankName,
// 		)

// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func LinkTransactionType(s database.Service, transactionID string, transactionTypeID string) error {
// 	query := `UPDATE transactions SET transaction_type_id = ? WHERE id = ?;`

// 	_, err := s.Exec(query, transactionTypeID, transactionID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil

// }

// func LinkBudgetToTransaction(s database.Service, transactionID string, budgetID string) error {
// 	query := `UPDATE transactions SET budget_item_id = ? WHERE id = ?;`

// 	_, err := s.Exec(query, budgetID, transactionID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
