package models

type TransactionType struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
}
