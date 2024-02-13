package models

import (
	"time"
)

type Transaction struct {
	ID                int       `json:"id"`
	TransactionType   string    `json:"transactionType"`
	TransactionDate   time.Time `json:"transactionDate"`
	TransactionAmount float64   `json:"transactionAmount"`
	TransactionID     string    `json:"transactionID"`
	TransactionName   string    `json:"transactionName"`
	TransactionMemo   string    `json:"transactionMemo"`
	CreatedAt         time.Time `json:"createdAt"`
	TransactionTypeID int       `json:"transactionTypeID"`
	BankName          string    `json:"bankName"`
}
