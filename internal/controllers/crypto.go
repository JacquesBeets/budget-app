package controllers

import (
	"budget-app/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ReturnAllCoinsView(ge *GinEngine, c *gin.Context) {
	var coins models.CryptoCoins
	r := ge.Router
	r.LoadHTMLFiles(Crypto)

	_, err := coins.FetchAll(ge.db())
	if err != nil {
		ge.ReturnErrorPage(c, err)
		return
	}

	c.HTML(http.StatusOK, "crypto_portfolio.html", gin.H{
		"now":         time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"CryptoCoins": coins,
	})
}

func (ge *GinEngine) ReturnCryptoView(c *gin.Context) {
	ReturnAllCoinsView(ge, c)
}

func (ge *GinEngine) SaveCryptoCoin(c *gin.Context) {
	var crypto *models.CryptoCoin

	coinID := c.PostForm("geckoid")
	coinName := c.PostForm("cryptoname")
	coinSymbol := c.PostForm("cryptosymbol")
	coinAmount := c.PostForm("amount")

	// Convert amount to float64
	amount, err := strconv.ParseFloat(coinAmount, 64)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	newCoin, err := crypto.New(coinID, coinName, coinSymbol, amount)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	_, err = newCoin.Create(ge.db())
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	ReturnAllCoinsView(ge, c)
}
