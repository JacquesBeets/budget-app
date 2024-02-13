package database

import (
	"fmt"
	"log"
)

func CreateBudgetTransactionsTable(s service) {
	// _, err := s.db.Exec(`
	//     DROP TABLE IF EXISTS budget_transactions;
	// `)

	// if err != nil {
	// 	log.Fatalf(fmt.Sprintf("Error dropping budget_transactions table: %v", err))
	// }

	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS budget_transactions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            transaction_id INTEGER NOT NULL,
            budget_id INTEGER NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating budget_transactions table: %v", err))
	}

}
