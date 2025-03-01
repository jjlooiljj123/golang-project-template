package mysql_test

import (
	"boilerplate/app/infrastructure/repositories/mysql"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockDriver implements driver.Driver and driver.Pinger for testing
type mockDriver struct {
	shouldFailOpen bool
	shouldFailPing bool
}

// mockConn implements driver.Conn and driver.Pinger
type mockConn struct {
	shouldFailPing bool
}

// Implement driver.Driver
func (m *mockDriver) Open(name string) (driver.Conn, error) {
	if m.shouldFailOpen {
		return nil, errors.New("mock open error")
	}
	return &mockConn{shouldFailPing: m.shouldFailPing}, nil
}

// Implement driver.Conn required methods
func (m *mockConn) Prepare(query string) (driver.Stmt, error) { return nil, nil }
func (m *mockConn) Close() error                              { return nil }
func (m *mockConn) Begin() (driver.Tx, error)                 { return nil, nil }

// Implement Ping for mockConn
func (m *mockConn) Ping(ctx context.Context) error {
	if m.shouldFailPing {
		return errors.New("mock ping error")
	}
	return nil
}

// mockStmt for completeness (minimal implementation)
type mockStmt struct{}

func (s *mockStmt) Close() error                                    { return nil }
func (s *mockStmt) NumInput() int                                   { return 0 }
func (s *mockStmt) Exec(args []driver.Value) (driver.Result, error) { return nil, nil }
func (s *mockStmt) Query(args []driver.Value) (driver.Rows, error)  { return nil, nil }

func TestOpenMySQLConnection(t *testing.T) {
	tests := []struct {
		name           string
		connString     string
		shouldFailOpen bool
		shouldFailPing bool
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Successful connection",
			connString:     "test:mysql@/dbname",
			shouldFailOpen: false,
			shouldFailPing: false,
			expectError:    false,
		},
		{
			name:           "Failed to open connection",
			connString:     "test:mysql@/dbname",
			shouldFailOpen: true,
			shouldFailPing: false,
			expectError:    true,
			expectedErrMsg: "failed to connect to database: mock open error",
		},
		{
			name:           "Failed to ping database",
			connString:     "test:mysql@/dbname",
			shouldFailOpen: false,
			shouldFailPing: true,
			expectError:    true,
			expectedErrMsg: "failed to ping database: mock ping error",
		},
		{
			name:           "Empty connection string",
			connString:     "",
			shouldFailOpen: false,
			shouldFailPing: false,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Register a unique mock driver for each test
			mockDriverName := fmt.Sprintf("mockmysql_%s", tt.name)
			sql.Register(mockDriverName, &mockDriver{
				shouldFailOpen: tt.shouldFailOpen,
				shouldFailPing: tt.shouldFailPing,
			})

			// Use our mock driver explicitly in the connection string
			testConnString := fmt.Sprintf("%s:%s", mockDriverName, tt.connString)

			// Call the function
			db, err := mysql.OpenMySQLConnection(testConnString)

			if tt.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErrMsg)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				// Clean up
				if db != nil {
					db.Close()
				}
			}
		})
	}
}
