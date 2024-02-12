package database

import (
	"budget-app/internal/models"
	"fmt"
	"log"
	"time"
)

func CreateBudgetTable(s service) {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS budget (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            amount FLOAT NOT NULL,
            created_at DATETIME NOT NULL,
            transaction_type_id INTEGER
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating transactions table: %v", err))
	}
}

func (s *service) GetBudget() ([]models.Budget, error) {
    rows, err := s.db.Query(`
        SELECT id, name, amount, created_at, transaction_type_id
        FROM budget;
    `)
    if err != nil {
        log.Fatalf(fmt.Sprintf("Error getting budget: %v", err))
        return nil, err
    }
    defer rows.Close()

    var budget []models.Budget
    for rows.Next() {
        var b models.Budget
        if err := rows.Scan(&b.ID, &b.Name, &b.Amount, &b.CreatedAt, &b.TransactionTypeID); err != nil {
            log.Fatalf(fmt.Sprintf("Error scanning budget: %v", err))
            return nil, err
        }
        budget = append(budget, b)
    }
    return budget, nil
}


func (s *service) SaveBudgetItem(budget models.Budget) error {
	_, err := s.db.Exec(`
        INSERT INTO budget (name, amount, created_at, transaction_type_id)
        VALUES (?, ?, ?, ?);
    `, budget.Name, budget.Amount, time.Now(), budget.TransactionTypeID)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error saving budget item: %v", err))
		return err
	}
	return nil
}
