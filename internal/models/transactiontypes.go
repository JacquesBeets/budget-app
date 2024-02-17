package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType struct {
	gorm.Model
	ID           uint          `gorm:"primaryKey" json:"id"`
	Title        string        `json:"title"`
	Category     *string       `json:"category"`
	Transactions []Transaction `gorm:"foreignKey:TransactionTypeID" json:"transactions"`
	CreatedAt    time.Time     // Automatically managed by GORM for creation time
	UpdatedAt    time.Time     // Automatically managed by GORM for creation time
}
