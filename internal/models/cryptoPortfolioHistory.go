package models

import "gorm.io/gorm"

type CryptoPortfolioHistory struct {
	gorm.Model
	CryptoCoinID                         uint       `json:"cryptoCoinID" gorm:"not null"`
	CryptoCoin                           CryptoCoin `json:"cryptoCoin" gorm:"foreignKey:CryptoCoinID"`
	CryptoCoinPrice                      float64    `json:"cryptoCoinPrice"`
	CryptoCoinAmount                     float64    `json:"cryptoCoinAmount"`
	CryptoCoinValue                      float64    `json:"cryptoCoinValue"`
	CryptoCoinValueZar                   float64    `json:"cryptoCoinValueZar"`
	CryptoCoinValueZarTotal              float64    `json:"cryptoCoinValueZarTotal"`
	CryptoCoinValueZarTotalChange        float64    `json:"cryptoCoinValueZarTotalChange"`
	CryptoCoinValueZarTotalChangePercent float64    `json:"cryptoCoinValueZarTotalChangePercent"`
}

type CryptoPortfolioHistories []CryptoPortfolioHistory

func (c *CryptoPortfolioHistory) New(
	coinID uint,
	coinPrice float64,
	coinAmount float64,
	coinValue float64,
	coinValueZar float64,
	coinValueZarTotal float64,
	coinValueZarTotalChange float64,
	coinValueZarTotalChangePercent float64,
) (*CryptoPortfolioHistory, error) {
	return &CryptoPortfolioHistory{
		CryptoCoinID:                         coinID,
		CryptoCoinPrice:                      coinPrice,
		CryptoCoinAmount:                     coinAmount,
		CryptoCoinValue:                      coinValue,
		CryptoCoinValueZar:                   coinValueZar,
		CryptoCoinValueZarTotal:              coinValueZarTotal,
		CryptoCoinValueZarTotalChange:        coinValueZarTotalChange,
		CryptoCoinValueZarTotalChangePercent: coinValueZarTotalChangePercent,
	}, nil
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
