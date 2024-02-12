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

func (t *Transaction) Update() error {
	return nil
}

func (t *Transaction) Delete() error {
	return nil
}

// func (s *database.Service) GetTransactions() ([]Transaction, error) {
// 	var transactions []Transaction
// 	s.db.Find(&transactions)
// 	return transactions, nil
// }
