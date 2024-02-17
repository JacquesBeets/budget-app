package models

import (
	"gorm.io/gorm"
)

type Budget struct {
	gorm.Model
	Name         string        `json:"name"`
	Amount       float64       `json:"amount"`
	Transactions []Transaction `gorm:"foreignKey:budget_id" json:"transactions"`
}
