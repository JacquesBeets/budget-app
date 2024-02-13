package models

import "time"

type Account struct {
	ID          int       `json:"id"`
	AccountName string    `json:"accountName"`
	AccountType string    `json:"accountType"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time `json:"createdAt"`
}
