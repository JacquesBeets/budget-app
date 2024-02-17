package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey" json:"id"`
	AccountName string    `json:"accountName"`
	AccountType string    `json:"accountType"`
	Balance     float64   `json:"balance"`
	CreatedAt   time.Time // Automatically managed by GORM for creation time
	UpdatedAt   time.Time // Automatically managed by GORM for creation time
}
