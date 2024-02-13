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
	SaveBudgetItem(budget models.Budget) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	GetDBPool() *sql.DB
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
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	s := &service{db: db}
	s.CreateTables()
	// Returning the service directly (which satisfies the Service interface)
	return s
}

func (s *service) GetDBPool() *sql.DB {
	return s.db
}

func (s *service) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.Query(query, args...)
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
	CreateAccountsTable(service{db: s.db})
	CreateBudgetTransactionsTable(service{db: s.db})
}
