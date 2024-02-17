package server

import (
	"budget-app/internal/controllers"
	"budget-app/internal/utils"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

var GinEngineVar *gin.Engine

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	GinEngineVar = r

	funcMap := template.FuncMap{
		"formatDate":                    utils.FormatDate,
		"formatPrice":                   utils.FormatPrice,
		"isEmpty":                       utils.IsEmpty,
		"isNil":                         utils.IsNil,
		"isTotalSpendGreaterThanBudget": utils.IsTotalSpendGreaterThanBudget,
		"dereferencePntr":               utils.DereferenceUintPtr,
	}
	r.SetFuncMap(funcMap)

	// handle static files
	r.Static("/static", "./static")

	// handle views
	HandleViews(r)
	HandleComponents(r)

	r.GET("/health", s.healthHandler)

	return GinEngineVar
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func HandleViews(r *gin.Engine) {
	views := controllers.ReturnViewsRouter(r)
	// authenticated := r.Group("/")
	// authenticated.GET("/transactions", controllers.GetTransactions)
	// authenticated.POST("/transactions", controllers.HandleOFXUpload)
	r.GET("/", views.HomePage)

	r.GET("templates/upload", views.UploadPage)
	r.GET("/upload", views.UploadPage)

	r.GET("templates/dashboard", views.Dashboard)
	r.GET("templates/transactions", views.ReturnTransactions)
	r.GET("templates/transactionstypes", views.ReturnTransactionTypes)

}

func HandleComponents(r *gin.Engine) {
	views := controllers.ReturnViewsRouter(r)

	r.GET("components/transactions", views.HandleTransctions)
	r.GET("components/budget/form", views.ReturnBudgetForm)

	// Posts
	r.POST("/upload/ofx", views.HandleOFXUpload)
	r.POST("/budget/add", views.SaveBudgetItem)
	r.POST("/transactionstypes/add", views.HandleTransactionTypeCreate)
	r.POST("/transactions/:id/transactionstypes/add", views.TransactionsAddTransactionType)
	r.POST("/transactions/:id/budgetitems/add", views.BudgetTransactionAdd)
	r.POST("/transactions/duplicate/remove", views.RemoveDuplicateTransactions)
}
