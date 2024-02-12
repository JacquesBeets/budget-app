package database

import (
	"budget-app/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

type Service interface {
	Health() map[string]string
	SaveMultipleTransactions(transactions []models.Transaction) error
	GetLatestTransactions() ([]models.Transaction, error)
	SaveBudgetItem(budget models.Budget) error
	GetBudget() ([]models.Budget, error)
}

type service struct {
	db *sql.DB
}

var (
	dburl = os.Getenv("DB_URL")
)

func New() Service {
	db, err := sql.Open("sqlite3", dburl)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	s := &service{db: db}
	s.CreateTables()
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) CreateTables() {
	CreateUserTable(service{db: s.db})
	CreateTransactionTable(service{db: s.db})
	CreateTransactionTypesTable(service{db: s.db})
	CreateBudgetTable(service{db: s.db})
}
