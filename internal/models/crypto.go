package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
)

type CryptoCoin struct {
	gorm.Model
	CoinGeckoCryptoID   string   `json:"cryptoID"`
	CryptoName          string   `json:"cryptoName"`
	CryptoSymbol        *string  `json:"cryptoSymbol"`
	CryptoAmountHolding *float64 `json:"cryptoAmount"`
	CryptoPrice         *float64 `json:"cryptoPrice"`
	CryptoPriceZar      *float64 `json:"cryptoPriceZar"`
	CurrentValueZar     *float64 `json:"currentValueZar"`
}

type CryptoCoins []CryptoCoin
type CryptoCoinIDs []uint

var (
	coinGeckoURL    = os.Getenv("COINGECKO_API_URL")
	coinGeckoAPIKey = os.Getenv("COINGECKO_API_KEY")
)

func (c *CryptoCoin) New(
	cryptoID string,
	cryptoName string,
	cryptoSymbol string,
	cryptoAmount float64,
) (*CryptoCoin, error) {
	return &CryptoCoin{
		CoinGeckoCryptoID:   cryptoID,
		CryptoName:          cryptoName,
		CryptoSymbol:        &cryptoSymbol,
		CryptoAmountHolding: &cryptoAmount,
	}, nil
}

func (c *CryptoCoins) FetchAll(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.Find(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoCoins) ReturnAllIds(tx *gorm.DB) (CryptoCoinIDs, error) {
	var ids CryptoCoinIDs
	response := tx.Model(&c).Select("id").Find(&ids)
	if response.Error != nil {
		return nil, response.Error
	}
	return ids, nil
}

func (c *CryptoCoins) Print() {
	for _, coin := range *c {
		println(coin.CryptoName, "===", coin.CoinGeckoCryptoID)
	}
}

func (c *CryptoCoins) UpdatePrices(tx *gorm.DB) (*gorm.DB, error) {
	var ids []string
	for _, coin := range *c {
		ids = append(ids, coin.CoinGeckoCryptoID)
	}
	stringofIDs := strings.Join(ids, ",")
	vs_currencies := "usd,zar"
	url := coinGeckoURL + "/simple/price?ids=" + stringofIDs + "&vs_currencies=" + vs_currencies + "&api_key=" + coinGeckoAPIKey + "&precision=18"
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return tx, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return tx, err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return tx, err
	}

	var result map[string]map[string]float64
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return tx, err
	}

	for coin, prices := range result {
		fmt.Printf("The price of %s in USD is %f\n", coin, prices["zar"])
		var coinModel CryptoCoin
		response := tx.Model(&coinModel).Where("coin_gecko_crypto_id", &coin).Updates(map[string]interface{}{
			"crypto_price":     prices["usd"],
			"crypto_price_zar": prices["zar"],
		})
		if response.Error != nil {
			return response, response.Error
		}
	}

	return tx, nil
}

func (c *CryptoCoin) Create(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.Create(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoCoin) FetchByID(tx *gorm.DB, id uint) *gorm.DB {
	return tx.First(&c, id)
}

func (c *CryptoCoin) Update(tx *gorm.DB, id uint) *gorm.DB {
	return tx.Save(c)
}

func (c *CryptoCoin) Delete(tx *gorm.DB, id uint) *gorm.DB {
	return tx.Delete(c, id)
}
