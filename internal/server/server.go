package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"budget-app/internal/database"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	newDB := database.New()
	NewServer := &Server{
		port: port,
		db:   *newDB, // Fix: Pass the value of newDB instead of its pointer
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
