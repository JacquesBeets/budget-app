package models


type BudgetTransaction struct {
    ID                int       `json:"id"`
    TransactionID     int       `json:"transactionID"`
    BudgetID          int       `json:"budgetID"`
}