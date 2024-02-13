package controllers

import (
	"budget-app/internal/database"

	"github.com/gin-gonic/gin"
)

type GinEngine struct {
	Router *gin.Engine
}

type Controller struct {
	db *database.Service
}

func NewController(db *database.Service) *Controller {
	return &Controller{
		db: db,
	}
}

func ReturnViewsRouter(r *gin.Engine) *GinEngine {
	return &GinEngine{Router: r}
}
