package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Admin    *bool  `gorm:"default:false" json:"admin"`
}
