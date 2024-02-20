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

	db, err := gorm.Open(sqlite.Open(dburl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	s := &Service{db: db}

	ServiceDB = db
	s.RunMigrations()

	return s
}

func ReturnDB() *gorm.DB {
	if ServiceDB == nil {
		log.Println("ReturnDB() is returning nil")
	}
	return ServiceDB
}

func (s *Service) Health() map[string]string {
	err := s.db.Exec("SELECT 1").Error
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *Service) RunMigrations() {
	// err := s.db.Migrator().DropTable(&models.Transactions{})
	// if err != nil {
	// 	log.Println(err)
	// }
	models := models.RegisteredModels
	s.db.AutoMigrate(models...)
}
