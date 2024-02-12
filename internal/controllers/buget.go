package controllers

import "budget-app/internal/models"

func NewBudget(
	name string,
	amount float64,
	transactionTypeID int,
) (*models.Budget, error) {

	return &models.Budget{
		Name:              name,
		Amount:            amount,
		TransactionTypeID: transactionTypeID,
	}, nil
}
