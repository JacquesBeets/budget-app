package controllers

import (
	"budget-app/internal/database"
	"budget-app/internal/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// func GetBudgetItems(s database.Service) ([]models.Budget, error) {
// 	db := s.GetDBPool()

// 	// rows, err := db.Query(`
// 	// 	SELECT id, name, amount, created_at, transaction_type_id
// 	// 	FROM budget;
// 	// `)

// 	stringQ := `SELECT b.id, b.name, b.amount, b.created_at, t.id, t.transaction_type, t.transaction_date, t.transaction_amount, t.transaction_id, t.transaction_name, t.transaction_memo, t.created_at, t.bank_name, t.transaction_type_id FROM budget b LEFT JOIN budget_transactions bt ON b.id = bt.budget_id LEFT JOIN transactions t ON bt.transaction_id = t.id`

// 	rows, err := db.Query(stringQ)

// 	if err != nil {
// 		fmt.Print("Error getting budget:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var budgets []models.Budget
// 	// for rows.Next() {
// 	// 	var b models.Budget
// 	// 	if err := rows.Scan(&b.ID, &b.Name, &b.Amount, &b.CreatedAt, &b.TransactionTypeID); err != nil {
// 	// 		log.Fatalf(fmt.Sprintf("Error scanning budget: %v", err))
// 	// 		return nil, err
// 	// 	}
// 	// 	budgets = append(budget, b)
// 	// }
// 	var currentBudgetID int
// 	var currentBudget *models.Budget
// 	var budgetName string
// 	var budgetAmount float64
// 	var budgetCreatedAt time.Time

// 	for rows.Next() {
// 		var budgetID int
// 		var transaction models.Transaction
// 		err := rows.Scan(&budgetID, &budgetName, &budgetAmount, &budgetCreatedAt, &transaction.ID, &transaction.TransactionType, &transaction.TransactionDate, &transaction.TransactionAmount, &transaction.TransactionID, &transaction.TransactionName, &transaction.TransactionMemo, &transaction.CreatedAt, &transaction.BankName, &transaction.TransactionTypeID)

// 		if err != nil {
// 			log.Fatalf(fmt.Sprintf("Error scanning budget: %v", err))
// 			return nil, err
// 		}

// 		fmt.Println("budgetID: ", budgetID, "currentBudgetID: ", currentBudgetID, "budgetName: ", budgetName, "budgetAmount: ", budgetAmount, "budgetCreatedAt: ", budgetCreatedAt, "transaction.ID: ", transaction.ID, "transaction.TransactionType: ", transaction.TransactionType, "transaction.TransactionDate: ", transaction.TransactionDate, "transaction.TransactionAmount: ", transaction.TransactionAmount, "transaction.TransactionID: ", transaction.TransactionID, "transaction.TransactionName: ", transaction.TransactionName, "transaction.TransactionMemo: ", transaction.TransactionMemo, "transaction.CreatedAt: ", transaction.CreatedAt, "transaction.BankName: ", transaction.BankName, "transaction.TransactionTypeID: ", transaction.TransactionTypeID)

// 		if currentBudgetID != budgetID {
// 			currentBudget = &models.Budget{
// 				ID:           budgetID,
// 				Name:         budgetName,
// 				Amount:       budgetAmount,
// 				CreatedAt:    budgetCreatedAt,
// 				Transactions: []models.Transaction{},
// 			}
// 			budgets = append(budgets, *currentBudget)
// 			currentBudgetID = budgetID
// 		}
// 		if transaction.ID.Valid {
// 			currentBudget.Transactions = append(currentBudget.Transactions, transaction)
// 		}
// 	}
// 	return budgets, nil
// }

// func NewBudget(
// 	name string,
// 	amount float64,
// 	transactionTypeID int,
// ) (*models.Budget, error) {

// 	return &models.Budget{
// 		Name:              name,
// 		Amount:            amount,
// 		TransactionTypeID: transactionTypeID,
// 		CreatedAt:         time.Now().UTC(),
// 	}, nil
// }

func (ge *GinEngine) SaveBudgetItem(c *gin.Context) {
	dbService := database.ReturnDB()
	var budget *models.Budget

	name := c.PostForm("name")
	amountStr := c.PostForm("amount")

	// Convert amount to float64
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Println("Error converting amount to float64: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount"})
		return
	}

	budget = &models.Budget{
		Name:   name,
		Amount: amount,
	}

	response := dbService.Create(budget).Scan(&budget)
	if response.Error != nil {
		fmt.Println("Error saving budget: ", response.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save budget"})
		return
	}

	c.JSON(http.StatusOK, budget)
}
