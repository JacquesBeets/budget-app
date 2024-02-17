package models

import (
	"gorm.io/gorm"
)

type BudgetTransaction struct {
	gorm.Model
	TransactionID int `json:"transactionID"`
	BudgetID      int `json:"budgetID"`
}
