package models

import "gorm.io/gorm"

type CryptoCoin struct {
	gorm.Model
	CoinGeckoCryptoID   string   `json:"cryptoID"`
	CryptoName          string   `json:"cryptoName"`
	CryptoSymbol        *string  `json:"cryptoSymbol"`
	CryptoAmountHolding *float64 `json:"cryptoAmount"`
	CryptoPrice         *float64 `json:"cryptoPrice"`
}

type CryptoCoins []CryptoCoin

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

func (c *CryptoCoin) Create(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.Create(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoCoin) FetchByID(tx *gorm.DB, id uint) *gorm.DB {
	return tx.Where("id = ?", id)
}

func (c *CryptoCoin) Update(tx *gorm.DB, id uint) *gorm.DB {
	return tx.Save(c)
}

func (c *CryptoCoin) Delete(tx *gorm.DB, id uint) *gorm.DB {
	return tx.Delete(c, id)
}
