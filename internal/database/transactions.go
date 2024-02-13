package database

import (
	"budget-app/internal/models"
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

func (s *service) SaveMultipleTransactions(transactions []models.Transaction) error {
	// Prepare the query for adding transactions
	query := `insert into transactions (
        transaction_type,
        transaction_date,
        transaction_amount,
        transaction_id,
        transaction_name,
        transaction_memo,
        created_at,
        transaction_type_id,
        bank_name
    ) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, transaction := range transactions {

		// Check if the transaction already exists
		existsQuery := `select count(*) from transactions where transaction_id = ?`
		var count int
		err = s.db.QueryRow(existsQuery, transaction.TransactionID).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			// Transaction already exists, do not add it
			log.Println("Transaction already exists")
			continue
		}

		// Prepare the query for adding transactions
		query := `insert into transactions (
			transaction_type,
			transaction_date,
			transaction_amount,
			transaction_id,
			transaction_name,
			transaction_memo,
			created_at,
			transaction_type_id,
			bank_name
		) values (?, ?, ?, ?, ?, ?, ?, ?, ?)`

		stmt, err := s.db.Prepare(query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			transaction.TransactionType,
			transaction.TransactionDate,
			transaction.TransactionAmount,
			transaction.TransactionID,
			transaction.TransactionName,
			transaction.TransactionMemo,
			transaction.CreatedAt,
			transaction.TransactionTypeID,
			transaction.BankName,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
