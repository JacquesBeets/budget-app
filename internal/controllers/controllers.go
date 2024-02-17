package controllers

import (
	"budget-app/internal/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GinEngine struct {
	Router *gin.Engine
}

func ReturnViewsRouter(r *gin.Engine) *GinEngine {
	return &GinEngine{Router: r}
}

func (ge *GinEngine) db() *gorm.DB {
	return database.ReturnDB()
}
