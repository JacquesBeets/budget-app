package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Views
const (
	IndexHTML        = "./views/index.html"
	SidenavHTML      = "./views/navigation/sidenav.html"
	BodyImportsHTML  = "./views/imports/bodyimports.html"
	HeadImportsHTML  = "./views/imports/headimports.html"
	DashboardHTML    = "./views/dashboard.html"
	UploadHTML       = "./views/uploads/upload.html"
	ErrorHTML        = "./views/error.html"
	Transactions     = "./views/transactions.html"
	TransactionTypes = "./views/transaction_types.html"
)

// Components
const (
	RecentTransactionComponent = "./views/components/recenttransactions.html"
	BudgetForm                 = "./views/components/budgetform.html"
)

const (
	StartDayOfMonth = "+21 days"
	DateNow         = "now"
	StringQuery     = "date(transaction_date) >= date(?, 'start of month', '-1 month', ?) AND date(transaction_date) <= date(?, 'start of month', ?)"
)

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

	r.LoadHTMLFiles(DashboardHTML)

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Now": time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
	})
}

func (ge *GinEngine) UploadPage(c *gin.Context) {
	r := ge.Router

	r.LoadHTMLFiles(UploadHTML)

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
