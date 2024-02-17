package database

import (
	"budget-app/internal/models"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

var (
	dburl     = os.Getenv("DB_URL")
	ServiceDB *gorm.DB
)

func New() *Service {
	// db, err := sql.Open("sqlite3", dburl)
	log.Println("New() is being called")
	db, err := gorm.Open(sqlite.Open(dburl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	s := &Service{db: db}
	// s.CreateTables()
	ServiceDB = db
	s.RunMigrations()
	log.Println("New() completed successfully")
	return s
}

func ReturnDB() *gorm.DB {
	if ServiceDB == nil {
		log.Println("ReturnDB() is returning nil")
	}
	return ServiceDB
}

// func (s *service) Query(query string, args ...interface{}) (*sql.Rows, error) {
// 	return s.db.Query(query, args...)
// }

// func (s *service) Exec(query string, args ...interface{}) (sql.Result, error) {
// 	return s.db.Exec(query, args...)
// }

func (s *Service) Health() map[string]string {
	err := s.db.Exec("SELECT 1").Error
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

// func (s *Service) CreateTables() {
// 	CreateUserTable(service{db: s.db})
// 	CreateTransactionTable(service{db: s.db})
// 	CreateTransactionTypesTable(service{db: s.db})
// 	CreateBudgetTable(service{db: s.db})
// 	CreateAccountsTable(service{db: s.db})
// 	CreateBudgetTransactionsTable(service{db: s.db})
// }

func (s *Service) RunMigrations() {
	models := models.RegisteredModels
	s.db.AutoMigrate(models...)
}
