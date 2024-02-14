package database

import (
	"budget-app/internal/models"
	"fmt"
	"log"
	"time"
)

func CreateAccountsTable(s service) {
	_, err := s.db.Exec(`
        DROP TABLE IF EXISTS accounts;
    `)

	if err != nil {
		log.Fatalf(fmt.Sprintf("Error dropping accounts table: %v", err))
	}

	_, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS accounts (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            account_name TEXT NOT NULL,
            account_type TEXT NOT NULL,
            balance FLOAT NOT NULL,
            created_at DATETIME NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating accounts table: %v", err))
	}

	SeedAccountsTable(s)

}

func SeedAccountsTable(s service) {
	currentAccounts := []models.Account{
		{
			AccountName: "FNB Checque",
			AccountType: "checking",
			Balance:     16400.89,
			CreatedAt:   time.Now(),
		},
		{
			AccountName: "FNB Savings",
			AccountType: "savings",
			Balance:     534.00,
			CreatedAt:   time.Now(),
		},
		{
			AccountName: "NEDBANK Checque",
			AccountType: "checking",
			Balance:     5284.71,
			CreatedAt:   time.Now(),
		},
		{
			AccountName: "NEDBANK Savings",
			AccountType: "savings",
			Balance:     204913.04,
			CreatedAt:   time.Now(),
		},
	}

	for _, account := range currentAccounts {
		_, err := s.db.Exec(`
            INSERT INTO accounts (account_name, account_type, balance, created_at)
            VALUES (?, ?, ?, ?);
        `, account.AccountName, account.AccountType, account.Balance, account.CreatedAt)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Error seeding accounts table: %v", err))
		}
	}

}
