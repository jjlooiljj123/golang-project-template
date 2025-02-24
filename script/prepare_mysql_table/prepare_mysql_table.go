package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var localhost string = "localhost"

// go run ./script/prepare_mysql_table/prepare_mysql_table.go
func main() {

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

	// masking the host to localhost (db is used in the .env for docker container connection in the docker compose)
	if dbHost == "db" {
		dbHost = localhost
	}

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

	// SQL statement to create a table
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS album (
        id VARCHAR(255) PRIMARY KEY,
        title VARCHAR(255) ,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`

	// Execute the SQL command
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Printf("Could not create table: %v", err)
	} else {
		fmt.Println("Table example_table created successfully")
	}
}
