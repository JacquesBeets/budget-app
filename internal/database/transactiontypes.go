package database

import (
	"fmt"
	"log"
)

func CreateTransactionTypesTable(s service) {
	// _, err := s.db.Exec(`
	// 	DROP TABLE IF EXISTS transactions_types;
	// `)
	// if err != nil {
	// 	log.Fatalf(fmt.Sprintf("Error dropping accounts table: %v", err))
	// }

	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS transactions_types (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
			category TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating transactions table: %v", err))
	}
}
