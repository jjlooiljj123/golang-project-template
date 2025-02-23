package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// go run ./script/insert_dummy_data/insert_dummy_data.go
func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve environment variables
	dbHost := os.Getenv("MYSQL_HOST")
	dbUser := os.Getenv("MYSQL_USER")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")

	// Connection string format: "username:password@tcp(host:port)/dbname"
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPassword, dbHost, dbName)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error on ping: ", err)
	}

	// Number of dummy records to insert
	numRecords := 10

	// Insert dummy data only if the record does not exist
	for i := 0; i < numRecords; i++ {
		// Generate dummy data
		id := fmt.Sprintf("A%04d", i+1) // Example: P0001, P0002, etc.
		title := fmt.Sprintf("Album Title %d", i+1)

		// SQL statement to insert a record if it does not already exist
		insertSQL := `
		INSERT INTO album (id, title)
		SELECT ?, ?
		WHERE NOT EXISTS (
				SELECT 1 FROM album WHERE id = ?
		);
		`

		// Execute the SQL command
		result, err := db.Exec(insertSQL, id, title, id)
		if err != nil {
			log.Printf("Could not check or insert record %d: %v", i+1, err)
		} else {
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected > 0 {
				fmt.Printf("Inserted new record: %s - %s\n", id, title)
			} else {
				fmt.Printf("Record already exists: %s - %s\n", id, title)
			}
		}
	}

	fmt.Printf("Attempted to insert %d dummy records into the albums table.\n", numRecords)
}
