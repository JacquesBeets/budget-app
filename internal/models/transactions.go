package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	ID                  uint             `gorm:"primaryKey" json:"id"`
	BankTransactionType *string          `json:"bankTransactionType"`
	TransactionDate     time.Time        `json:"transactionDate"`
	TransactionAmount   float64          `json:"transactionAmount"`
	BankTransactionID   string           `json:"transactionID"`
	TransactionName     *string          `json:"transactionName"`
	TransactionMemo     *string          `json:"transactionMemo"`
	BankName            string           `json:"bankName"`
	TransactionTypeID   *uint            `json:"transactionTypeID"`
	TransactionType     *TransactionType `json:"transactionType"`
	BudgetID            *uint            `json:"budgetID"` // Changed to uint
	Budget              Budget           `json:"budget"`
	AccountID           *uint            `json:"accountID"`
	CreatedAt           time.Time        // Automatically managed by GORM for creation time
	UpdatedAt           time.Time        // Automatically managed by GORM for creation time
}
