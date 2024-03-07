package controllers

import (
	"budget-app/internal/models"
	"budget-app/internal/utils"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CoinData struct {
	CoinHistory   models.CryptoPortfolioHistory
	Coin          models.CryptoCoin
	PercentChange float64
}

func ReturnAllCoinsView(ge *GinEngine, c *gin.Context) {
	var coins models.CryptoCoins
	var cryptoHistories models.CryptoPortfolioHistories
	r := ge.Router
	r.LoadHTMLFiles(Crypto)

	coinIds, err := coins.ReturnAllIds(ge.db())
	if err != nil {
		ge.ReturnErrorPage(c, err)
		return
	}

	_, err = cryptoHistories.FetchAll(ge.db(), coinIds)
	if err != nil {
		ge.ReturnErrorPage(c, err)
		return
	}

	totalValue := 0.0
	prevTotalValue := 0.0
	bestPerformingCoin := models.CryptoPortfolioHistory{}
	coinData := []CoinData{}
	for i := range cryptoHistories {
		coinHist := cryptoHistories[i]

		// Calculate the current value of the coin
		if coinHist.CryptoCoin.CryptoPriceZar != nil && coinHist.CryptoCoin.CryptoAmountHolding != nil {
			temp := *coinHist.CryptoCoin.CryptoPriceZar * *coinHist.CryptoCoin.CryptoAmountHolding
			coinHist.CryptoCoin.CurrentValueZar = &temp
		}
		totalValue += *coinHist.CryptoCoin.CurrentValueZar

		// Calculate the previous total value
		prevTotalValue += coinHist.CryptoCoinPriceZar * coinHist.CryptoCoinAmountHolding

		// Find the best performing coin
		if bestPerformingCoin.ID == 0 || coinHist.CalculatePercentageChange() > bestPerformingCoin.CalculatePercentageChange() {
			bestPerformingCoin = coinHist
		}

		coinData = append(coinData, CoinData{
			CoinHistory:   coinHist,
			Coin:          coinHist.CryptoCoin,
			PercentChange: coinHist.CalculatePercentageChange(),
		})
	}

	sort.Slice(coinData, func(i, j int) bool {
		// Check if the values are not nil to avoid a runtime panic
		if coinData[i].Coin.CurrentValueZar != nil && coinData[j].Coin.CurrentValueZar != nil {
			return *coinData[i].Coin.CurrentValueZar > *coinData[j].Coin.CurrentValueZar
		}
		return false
	})

	c.HTML(http.StatusOK, "crypto_portfolio.html", gin.H{
		"CryptoCoins":                 coins,
		"TotalValue":                  totalValue,
		"CoinData":                    coinData,
		"PrevTotalValue":              prevTotalValue,
		"TotalValuePercentChange":     utils.CalculatePercentageChange(prevTotalValue, totalValue),
		"PrevTotalValuePercentChange": utils.CalculatePercentageChange(totalValue, prevTotalValue),
		"BestPerformingCoin":          bestPerformingCoin,
		"BestPerformingCoinPrice":     *bestPerformingCoin.CryptoCoin.CryptoPriceZar * *bestPerformingCoin.CryptoCoin.CryptoAmountHolding,
		"BestPerformingCoinPercent":   bestPerformingCoin.CalculatePercentageChange(),
	})
}

func (ge *GinEngine) FetchCurrentCrypoPrices(c *gin.Context) {
	cryptoCoins := models.CryptoCoins{}
	cryptoHistories := models.CryptoPortfolioHistories{}

	_, err := cryptoCoins.FetchAll(ge.db())
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	_, err = cryptoCoins.UpdatePrices(ge.db())
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	newHistories, err := cryptoHistories.New(ge.db(), cryptoCoins)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
		return
	}

	_, err = newHistories.Save(ge.db())
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

func (ge *GinEngine) ReturnCryptoModal(c *gin.Context) {
	r := ge.Router
	r.LoadHTMLFiles(CryptoModal)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ge.ReturnErrorJSON(c, err)
	}
	

	if idStr == "new" {
		c.HTML(http.StatusOK, "cryptomodal.html", gin.H{})
		return
	}

	coin := models.CryptoCoin{}
	coin.ID = uint(id)

	response := coin.FetchByID(ge.db(), uint(id))
	if response.Error != nil {
		ge.ReturnErrorJSON(c, response.Error)
		return
	}

	c.HTML(http.StatusOK, "cryptomodal.html", gin.H{
		"Coin": coin,
	})

}
