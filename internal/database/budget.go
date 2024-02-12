package database

import (
	"fmt"
	"log"
)

func CreateBudgetTable(s service) {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS buget (
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
