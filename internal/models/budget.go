package models

import "time"

type Budget struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	Amount            float64   `json:"amount"`
	CreatedAt         time.Time `json:"createdAt"`
	TransactionTypeID int       `json:"transactionTypeID"`
}
