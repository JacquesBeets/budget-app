package controllers

import "github.com/gin-gonic/gin"

type GinEngine struct {
	Router *gin.Engine
}

func ReturnViewsRouter (r *gin.Engine) *GinEngine {
	return &GinEngine{Router: r}
}
