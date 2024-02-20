package models

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	BankTransactionType *string          `json:"bankTransactionType"`
	TransactionDate     time.Time        `json:"transactionDate"`
	TransactionAmount   float64          `json:"transactionAmount"`
	BankTransactionID   string           `json:"transactionID"`
	TransactionName     *string          `json:"transactionName"`
	TransactionMemo     *string          `json:"transactionMemo"`
	BankName            string           `json:"bankName"`
	TransactionTypeID   *uint            `json:"transactionTypeID"`
	TransactionType     *TransactionType `json:"transactionType"`
	BudgetID            *uint            `json:"budgetID"`
	Budget              *Budget          `json:"budget"`
	AccountID           *uint            `json:"accountID"`
	AutoCategorized     bool             `gorm:"default:false" json:"autoCategorized"`
}

type Transactions []Transaction

func (t *Transaction) New(
	transactionType string,
	transactionDate string,
	transactionAmount float64,
	transactionID string,
	transactionName string,
	transactionMemo string,
	bankName string,
) (*Transaction, error) {
	transactionDateParsed, err := time.Parse("2006-01-02", transactionDate)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		BankTransactionType: &transactionType,
		TransactionDate:     transactionDateParsed,
		TransactionAmount:   transactionAmount,
		BankTransactionID:   transactionID,
		TransactionName:     &transactionName,
		TransactionMemo:     &transactionMemo,
		BankName:            bankName,
	}, nil
}

func (t *Transaction) Exists(tx *gorm.DB) bool {
	var similarTransaction Transaction
	var response *gorm.DB

	if t.BankName == "FNB" {
		response = tx.Where("date(transaction_date) = date(?) AND transaction_amount = ? AND transaction_memo = ?", &t.TransactionDate, &t.TransactionAmount, &t.TransactionMemo).First(&similarTransaction)
	} else {
		response = tx.Where("date(transaction_date) = date(?) AND transaction_amount = ? AND transaction_name = ?", &t.TransactionDate, &t.TransactionAmount, &t.TransactionName).First(&similarTransaction)
	}

	if response.Error != nil {
		if errors.Is(response.Error, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}

	return true
}

func (t *Transaction) AutoCategorize(tx *gorm.DB) {
	var similarTransaction Transaction
	if t.BankName == "FNB" {
		response := tx.Where("transaction_memo = ?", t.TransactionMemo).First(&similarTransaction).Scan(&similarTransaction)
		if response.Error != nil {
			fmt.Printf("could not find similar transaction: %v", response.Error)
			return
		}

		if similarTransaction.BudgetID != nil {
			t.BudgetID = similarTransaction.BudgetID
			t.AutoCategorized = true
		}

		if similarTransaction.TransactionTypeID != nil {
			t.TransactionTypeID = similarTransaction.TransactionTypeID
			t.AutoCategorized = true
		}
	} else {
		response := tx.Where("transaction_name = ?", t.TransactionName).First(&similarTransaction).Scan(&similarTransaction)
		if response.Error != nil {
			fmt.Printf("could not find similar transaction: %v", response.Error)
			return
		}

		if similarTransaction.BudgetID != nil {
			t.BudgetID = similarTransaction.BudgetID
			t.AutoCategorized = true
		}

		if similarTransaction.TransactionTypeID != nil {
			t.TransactionTypeID = similarTransaction.TransactionTypeID
			t.AutoCategorized = true
		}
	}
}

func (t *Transaction) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&t).Scan(&t)
}

func (t *Transaction) Update(tx *gorm.DB) *gorm.DB {
	return tx.Save(&t).Scan(&t)
}

func (t *Transaction) FindByID(tx *gorm.DB, id uint) *gorm.DB {
	return tx.First(&t, &t.ID).Scan(&t)
}

func (t *Transactions) PrintAll() {
	fmt.Println("Printing all transactions")
	for i, transaction := range *t {
		fmt.Printf("Transaction %d: %+v\n", i+1, transaction)
	}
}

func (t *Transactions) Create(tx *gorm.DB) *gorm.DB {
	return tx.Create(&t).Scan(&t)
}
