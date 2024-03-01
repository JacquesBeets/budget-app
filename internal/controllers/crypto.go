package controllers

import (
	"budget-app/internal/models"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ReturnAllCoinsView(ge *GinEngine, c *gin.Context) {
	var coins models.CryptoCoins
	var cryptoHistory models.CryptoPortfolioHistories
	r := ge.Router
	r.LoadHTMLFiles(Crypto)

	_, err := coins.FetchAll(ge.db())
	if err != nil {
		ge.ReturnErrorPage(c, err)
		return
	}

	_, err = cryptoHistory.FetchAll(ge.db())
	if err != nil {
		ge.ReturnErrorPage(c, err)
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"OK":    http.StatusOK,
	// 	"VALUE": cryptoHistory,
	// })

	totalValue := 0.0
	// Calculate the current value of the coins
	for i := range coins {
		if coins[i].CryptoPriceZar != nil && coins[i].CryptoAmountHolding != nil {
			temp := *coins[i].CryptoPriceZar * *coins[i].CryptoAmountHolding
			coins[i].CurrentValueZar = &temp
		}
		totalValue += *coins[i].CurrentValueZar
	}

	sort.Slice(coins, func(i, j int) bool {
		// Check if the values are not nil to avoid a runtime panic
		if coins[i].CurrentValueZar != nil && coins[j].CurrentValueZar != nil {
			return *coins[i].CurrentValueZar > *coins[j].CurrentValueZar
		}
		return false
	})

	c.HTML(http.StatusOK, "crypto_portfolio.html", gin.H{
		"now":         time.Date(2017, 0o7, 0o1, 0, 0, 0, 0, time.UTC),
		"CryptoCoins": coins,
		"TotalValue":  totalValue,
	})
}

func (ge *GinEngine) FetchCurrentCrypoPrices(c *gin.Context) {
	cryptoCoins := models.CryptoCoins{}
	cryptoHistory := models.CryptoPortfolioHistories{}

	_, err := cryptoCoins.FetchAll(ge.db())
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	newHistories, err := cryptoHistory.New(ge.db(), cryptoCoins)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	_, err = newHistories.Save(ge.db())
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	_, err = cryptoCoins.UpdatePrices(ge.db())
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	ReturnAllCoinsView(ge, c)
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
