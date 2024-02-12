package controllers

import (
	"budget-app/internal/utils"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Views
const (
	IndexHTML       = "./views/index.html"
	SidenavHTML     = "./views/navigation/sidenav.html"
	BodyImportsHTML = "./views/imports/bodyimports.html"
	HeadImportsHTML = "./views/imports/headimports.html"
	DashboardHTML   = "./views/dashboard.html"
	UploadHTML      = "./views/uploads/upload.html"
	ErrorHTML       = "./views/error.html"
)

// Components
const (
	RecentTransactionComponent = "./views/components/recenttransactions.html"
	BudgetForm                 = "./views/components/budgetform.html"
)

func ParseFiles(files ...string) *template.Template {
	return template.Must(template.ParseFiles(files...))
}

type PageData struct {
	Version string
}

func (ge *GinEngine) HomePage(c *gin.Context) {

	r := ge.Router

	// Generate the current timestamp in milliseconds
	version := time.Now().UnixNano() / int64(time.Millisecond)

	r.LoadHTMLFiles(IndexHTML, SidenavHTML, BodyImportsHTML, HeadImportsHTML)

	c.HTML(http.StatusOK, "homepage/index.html", gin.H{
		"Version": fmt.Sprintf("%d", version),
	})
}

func (ge *GinEngine) Dashboard(c *gin.Context) {
	r := ge.Router

	tmpl := ParseFiles(DashboardHTML)

	r.SetHTMLTemplate(tmpl)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func ReturnUploadTemp() *template.Template {
	return ParseFiles(UploadHTML)
}

func (ge *GinEngine) UploadPage(c *gin.Context) {
	r := ge.Router

	tmpl := ReturnUploadTemp()

	r.SetHTMLTemplate(tmpl)

	c.HTML(http.StatusOK, "upload.html", gin.H{
		"now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func (ge *GinEngine) UploadPageRefreshed(c *gin.Context) {
	r := ge.Router

	// Generate the current timestamp in milliseconds
	version := time.Now().UnixNano() / int64(time.Millisecond)

	r.LoadHTMLFiles(IndexHTML, UploadHTML)

	c.HTML(http.StatusOK, "homepage/index.html", gin.H{
		"Version": fmt.Sprintf("%d", version),
	})
}

func (ge *GinEngine) HandleTransctions(c *gin.Context) {
	r := ge.Router
	funcMap := template.FuncMap{
		"formatDate":  utils.FormatDate,
		"formatPrice": utils.FormatPrice,
	}
	r.SetFuncMap(funcMap)
	r.LoadHTMLFiles(RecentTransactionComponent)

	transactions, err := GetTransactions()
	if err != nil {
		r.LoadHTMLFiles(ErrorHTML)
		c.HTML(http.StatusInternalServerError, "views/error.html", gin.H{"error": "could not fetch transactions"})
		return
	}

	recentTotal := 0.0
	for _, t := range transactions {
		recentTotal += float64(t.TransactionAmount)
	}

	c.HTML(http.StatusOK, "recenttransactions.html", gin.H{
		"now":              time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"RecentTotal":      recentTotal,
		"BudgetTotal":      "44300.00",
		"Transactions":     transactions,
		"TransactionCount": len(transactions),
	})
}

func (ge *GinEngine) ReturnBudgetForm(c *gin.Context) {
	// single file
	r := ge.Router
	r.LoadHTMLFiles(BudgetForm)

	c.HTML(http.StatusOK, "budgetform.html", gin.H{
		"Version": "1",
	})
}
