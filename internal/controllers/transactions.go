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
		TransactionType:   transactionType,
		TransactionDate:   transactionDateParsed,
		TransactionAmount: transactionAmount,
		TransactionID:     transactionID,
		TransactionName:   transactionName,
		TransactionMemo:   transactionMemo,
		BankName:          bankName,
		CreatedAt:         time.Now().UTC(),
	}, nil
}

func GetTransactions() ([]models.Transaction, error) {
	var dbService database.Service = database.New()
	fmt.Println("Getting transactions")
	trns, err := dbService.GetLatestTransactions()
	if err != nil {
		return nil, err
	}
	return trns, nil
}

func (ge *GinEngine) HandleOFXUpload(c *gin.Context) {
	// single file
	r := ge.Router
	file, _ := c.FormFile("ofx")
	bank := c.PostForm("bank")
	dst := "./uploads/" + bank + "-" + file.Filename

	// Upload the file to specific dst.
	c.SaveUploadedFile(file, dst)

	tmpl := ReturnUploadTemp()

	r.SetHTMLTemplate(tmpl)

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
	var dbService database.Service = database.New()

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
				if err := dbService.SaveMultipleTransactions(Transactions); err != nil {
					fmt.Printf("could not add transaction: %v", err)
					return err
				}
			}
		}
	}

	return nil
}
