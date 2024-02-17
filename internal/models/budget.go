package models

import (
	"time"

	"gorm.io/gorm"
)

type Budget struct {
	gorm.Model
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `json:"name"`
	Amount       float64       `json:"amount"`
	Transactions []Transaction `gorm:"foreignKey:BudgetID" json:"transactions"`
	CreatedAt    time.Time     // Automatically managed by GORM for creation time
	UpdatedAt    time.Time     // Automatically managed by GORM for creation time
}
