package models

import "gorm.io/gorm"

type CryptoPortfolioHistory struct {
	gorm.Model
	CryptoCoinID            uint       `json:"cryptoCoinID" gorm:"not null"`
	CryptoCoin              CryptoCoin `json:"cryptoCoin" gorm:"foreignKey:CryptoCoinID"`
	CryptoCoinPrice         float64    `json:"cryptoCoinPrice"`
	CryptoCoinPriceZar      float64    `json:"cryptoCoinPriceZar"`
	CryptoCoinAmountHolding float64    `json:"cryptoCoinAmountHolding"`
}

type CryptoPortfolioHistories []CryptoPortfolioHistory

func (c *CryptoPortfolioHistory) New(
	coinID uint,
	coinPrice float64,
	coinAmount float64,
	coinValueZar float64,
) (*CryptoPortfolioHistory, error) {
	return &CryptoPortfolioHistory{
		CryptoCoinID:                         coinID,
		CryptoCoinPrice:                      coinPrice,
		CryptoCoinPriceZar:                   coinValueZar,
		CryptoCoinAmountHolding:              coinAmount,
	}, nil
}

func (c *CryptoPortfolioHistories) New(tx *gorm.DB, freshCoins CryptoCoins) (*CryptoPortfolioHistories, error) {
	var cryptoPortfolioHistories CryptoPortfolioHistories
	for _, coin := range freshCoins {
		cryptoPortfolioHistory := CryptoPortfolioHistory{
			CryptoCoinID:                         coin.ID,
			CryptoCoinPrice:                      *coin.CryptoPrice,
			CryptoCoinPriceZar:                   *coin.CryptoPriceZar,
			CryptoCoinAmountHolding:              *coin.CryptoAmountHolding,
		}
		cryptoPortfolioHistories = append(cryptoPortfolioHistories, cryptoPortfolioHistory)
	}
	return &cryptoPortfolioHistories, nil
}

func (c *CryptoPortfolioHistories) FetchAll(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.Preload("CryptoCoin").Order("created_at desc").Find(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistory) FetchOne(tx *gorm.DB, id uint) (*gorm.DB, error) {
	response := tx.Preload("CryptoCoin").First(&c, id).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistory) Update(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.Model(&c).Updates(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistory) Delete(tx *gorm.DB, id uint) (*gorm.DB, error) {
	response := tx.Delete(&c, id)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistories) Save(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.Create(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistories) FetchAllByCoinID(tx *gorm.DB, coinID uint) (*gorm.DB, error) {
	response := tx.Preload("CryptoCoin").Where("crypto_coin_id = ?", coinID).Order("created_at desc").Find(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistories) FetchAllByCoinIDAndDate(tx *gorm.DB, coinID uint, date string) (*gorm.DB, error) {
	response := tx.Preload("CryptoCoin").Where("crypto_coin_id = ? AND date(created_at) = ?", coinID, date).Order("created_at desc").Find(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistories) FetchAllByDate(tx *gorm.DB, date string) (*gorm.DB, error) {
	response := tx.Preload("CryptoCoin").Where("date(created_at) = ?", date).Order("created_at desc").Find(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioHistories) FetchAllByDateRange(tx *gorm.DB, startDate string, endDate string) (*gorm.DB, error) {
	response := tx.Preload("CryptoCoin").Where("date(created_at) >= ? AND date(created_at) <= ?", startDate, endDate).Order("created_at desc").Find(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}
