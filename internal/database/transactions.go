package database

import (
	"fmt"
	"log"
)

func CreateTransactionTable(s service) {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS transactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            transaction_type,
            transaction_date,
            transaction_amount FLOAT NOT NULL,
            transaction_id TEXT,
            transaction_name TEXT,
            transaction_memo TEXT,
            created_at DATETIME NOT NULL,
            bank_name TEXT NOT NULL,
            transaction_type_id INTEGER,
            FOREIGN KEY (transaction_type_id) REFERENCES transaction_types(id) 
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating transactions table: %v", err))
	}
}
