package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// OpenMySQLConnection establishes a connection to the MySQL database
func OpenMySQLConnection(connectionString string) (*sql.DB, error) {
	fmt.Printf("connectionString: %s \n", connectionString)
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Ping the database to ensure the connection is good
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
