package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	BankTransactionType *string          `json:"bankTransactionType"`
	TransactionDate     time.Time        `json:"transactionDate"`
	TransactionAmount   float64          `json:"transactionAmount"`
	BankTransactionID   string           `gorm:"unique" json:"transactionID"`
	TransactionName     *string          `json:"transactionName"`
	TransactionMemo     *string          `json:"transactionMemo"`
	BankName            string           `json:"bankName"`
	TransactionTypeID   *uint            `json:"transactionTypeID"`
	TransactionType     *TransactionType `json:"transactionType"`
	BudgetID            *uint            `json:"budgetID"`
	Budget              *Budget          `json:"budget"`
	AccountID           *uint            `json:"accountID"`
}
