package server

import (
	"budget-app/internal/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

var GinEngineVar *gin.Engine

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	GinEngineVar = r
	// handle static files
	r.Static("/static", "./static")

	// handle views
	HandleViews(r)
	HandleComponents(r)

	r.GET("/health", s.healthHandler)

	// authenticated := r.Group("/")
	// authenticated.GET("/transactions", controllers.GetTransactions)
	// authenticated.POST("/transactions", controllers.HandleOFXUpload)
	return GinEngineVar
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func HandleViews(r *gin.Engine) {
	views := controllers.ReturnViewsRouter(r)
	r.GET("/", views.HomePage)

	r.GET("templates/upload", views.UploadPage)
	r.GET("/upload", views.UploadPage)

	r.GET("templates/dashboard", views.Dashboard)

	// Posts
	r.POST("/upload/ofx", views.HandleOFXUpload)

}

func HandleComponents(r *gin.Engine) {
	views := controllers.ReturnViewsRouter(r)

	r.GET("components/transactions", views.HandleTransctions)
}
