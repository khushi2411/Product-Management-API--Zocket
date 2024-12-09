package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	// Replace with your actual database credentials
	connStr := "host=localhost port=5432 user=postgres password=khushi2411 dbname=products_db sslmode=disable"

	// Initialize the database connection
	tempDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to the database (host=localhost port=5432): %v", err)
	}

	// Verify the connection
	err = tempDB.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database: %v", err)
	}

	// Assign to the global DB variable
	DB = tempDB
	fmt.Println("Successfully connected to the database!")
}

// CloseDB closes the database connection
func CloseDB() {
	if err := DB.Close(); err != nil {
		log.Fatalf("Error while closing the database connection: %v", err)
	}
	fmt.Println("Database connection closed successfully.")
}
