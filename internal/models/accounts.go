package models

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	AccountName string  `json:"accountName"`
	AccountType string  `json:"accountType"`
	Balance     float64 `json:"balance"`
}
