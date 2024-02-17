package models

import (
	"time"

	"gorm.io/gorm"
)

type BudgetTransaction struct {
	gorm.Model
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionID int       `json:"transactionID"`
	BudgetID      int       `json:"budgetID"`
	CreatedAt     time.Time // Automatically managed by GORM for creation time
	UpdatedAt     time.Time // Automatically managed by GORM for creation time
}
