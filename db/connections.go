package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// host=34.142.51.130
// host=localhost
func Connect() {
	var err error
	connStr := "user=scissorhands_user password=securepassword dbname=scissorhands host=34.142.51.130 sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	log.Println("Database connection successful")
}
