package models

import (
	"gorm.io/gorm"
)

type TransactionType struct {
	gorm.Model
	Title        string        `json:"title"`
	Category     *string       `json:"category"`
	Transactions []Transaction `gorm:"foreignKey:TransactionTypeID" json:"transactions"`
}
