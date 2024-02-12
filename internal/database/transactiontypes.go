package database

import (
	"fmt"
	"log"
)

func CreateTransactionTypesTable(s service) {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS transactions_types (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating transactions table: %v", err))
	}
}
