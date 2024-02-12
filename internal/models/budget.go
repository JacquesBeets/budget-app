package models

type Budget struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Amount            float64 `json:"amount"`
	CreatedAt         string  `json:"createdAt"`
	TransactionTypeID int     `json:"transactionTypeID"`
}
