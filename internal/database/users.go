package database

import (
	"fmt"
	"log"
)

func CreateUserTable(s service) {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            email TEXT NOT NULL UNIQUE,
            password TEXT NOT NULL,
            admin BOOLEAN NOT NULL DEFAULT FALSE
        );
    `)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error creating users table: %v", err))
	}
}
