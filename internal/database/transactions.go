package database

import (
	"budget-app/internal/models"
	"fmt"
	"log"
	"time"
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

func (s *service) GetLatestTransactions() ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	currentMonth := time.Now().Format("01") // "01" is the format for two-digit month in Go

	query := `SELECT * 
	FROM transactions 
	WHERE date(transaction_date) >= date('now', 'start of month', '-1 month', '+23 days') 
	ORDER BY date(transaction_date) DESC;`
	rows, err := s.db.Query(query, currentMonth)
	if err != nil {
		fmt.Println(err)
		return transactions, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction models.Transaction
		var transactionDateString string

		err := rows.Scan(
			&transaction.ID,
			&transaction.TransactionType,
			&transactionDateString,
			&transaction.TransactionAmount,
			&transaction.TransactionID,
			&transaction.TransactionName,
			&transaction.TransactionMemo,
			&transaction.CreatedAt,
			&transaction.BankName,
			&transaction.TransactionTypeID,
		)
		if err != nil {
			fmt.Println(err)
			return transactions, err
		}

		// Parse the date string into a time.Time type
		transaction.TransactionDate, err = time.Parse("2006-01-02 15:04:05-07:00", transactionDateString)
		if err != nil {
			fmt.Println(err)
			return transactions, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
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
