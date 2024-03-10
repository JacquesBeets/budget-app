package models

import (
	"gorm.io/gorm"
)

type CryptoPortfolioTotalHistory struct {
	gorm.Model
	TotalValue    float64 `json:"totalValue"`
	TotalValueZar float64 `json:"totalValueZar"`
}

type CryptoPortfolioTotalHistories []CryptoPortfolioTotalHistory

func NewCryptoPortfolioTotalHistory(
	totalValue float64,
	totalValueZar float64,
) *CryptoPortfolioTotalHistory {
	return &CryptoPortfolioTotalHistory{
		TotalValue:    totalValue,
		TotalValueZar: totalValueZar,
	}
}

func (c *CryptoPortfolioTotalHistory) FetchOne(tx *gorm.DB) (*gorm.DB, error) {
	response := tx.First(&c).Scan(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioTotalHistory) Create(tx *gorm.DB, pt *PortfolioTotal) (*gorm.DB, error) {
	c.TotalValue = pt.TotalValue
	c.TotalValueZar = pt.TotalValueZar
	response := tx.Create(&c)
	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}

func (c *CryptoPortfolioTotalHistories) Create(tx *gorm.DB, totalsMap *PortfolioTotalsMap) {
	var newHistories CryptoPortfolioTotalHistories

	for _, total := range *totalsMap {
		newTotal := NewCryptoPortfolioTotalHistory(total.TotalValue, total.TotalValueZar)
		newHistories = append(newHistories, *newTotal)
	}
	tx.Create(&newHistories)
}

func (c *CryptoPortfolioTotalHistories) FetchAll(tx *gorm.DB) (*gorm.DB, error) {

	response := tx.Order("id DESC").Limit(20).Find(&c)

	// Reverse the order of items in c
	for i, j := 0, len(*c)-1; i < j; i, j = i+1, j-1 {
		(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
	}

	if response.Error != nil {
		return response, response.Error
	}
	return response, nil
}
