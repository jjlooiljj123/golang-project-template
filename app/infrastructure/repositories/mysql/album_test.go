package mysql_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"boilerplate/app/domain/entity"
	customerr "boilerplate/app/domain/errors"
	"boilerplate/app/infrastructure/repositories/mysql"
)

func TestAlbumRepository(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		setupMock      func(sqlmock.Sqlmock)
		action         func(*mysql.AlbumRepository) interface{}
		expectedResult interface{}
		expectError    bool
		expectedErr    error
	}{
		// GetAlbums tests
		{
			name: "GetAlbums_Success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow("1", "Album1").
					AddRow("2", "Album2")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM album")).
					WillReturnRows(rows)
			},
			action: func(r *mysql.AlbumRepository) interface{} {
				albums, _ := r.GetAlbums(context.Background())
				return albums
			},
			expectedResult: []entity.Album{
				{ID: entity.AlbumID("1"), Title: "Album1"},
				{ID: entity.AlbumID("2"), Title: "Album2"},
			},
			expectError: false,
		},
		{
			name: "GetAlbums_QueryError",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM album")).
					WillReturnError(errors.New("query error"))
			},
			action: func(r *mysql.AlbumRepository) interface{} {
				_, err := r.GetAlbums(context.Background())
				return err
			},
			expectedResult: nil,
			expectError:    true,
			expectedErr:    fmt.Errorf("error querying data: query error"),
		},

		// CreateAlbum tests
		{
			name: "CreateAlbum_Success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO album (id, title) VALUES (?, ?)")).
					ExpectExec().
					WithArgs("1", "New Album").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			action: func(r *mysql.AlbumRepository) interface{} {
				id, _ := r.CreateAlbum(context.Background(), entity.Album{
					ID:    entity.AlbumID("1"),
					Title: "New Album",
				})
				return id
			},
			expectedResult: "1",
			expectError:    false,
		},
		{
			name: "CreateAlbum_ExecError",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO album (id, title) VALUES (?, ?)")).
					ExpectExec().
					WithArgs("1", "New Album").
					WillReturnError(errors.New("exec error"))
			},
			action: func(r *mysql.AlbumRepository) interface{} {
				_, err := r.CreateAlbum(context.Background(), entity.Album{
					ID:    entity.AlbumID("1"),
					Title: "New Album",
				})
				return err
			},
			expectedResult: nil,
			expectError:    true,
			expectedErr:    fmt.Errorf("error executing insert: exec error"),
		},

		// GetAlbumByID tests
		{
			name: "GetAlbumByID_Success",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title"}).
					AddRow("1", "Test Album")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM album WHERE id = ?")).
					WithArgs("1").
					WillReturnRows(rows)
			},
			action: func(r *mysql.AlbumRepository) interface{} {
				album, _ := r.GetAlbumByID(context.Background(), "1")
				return album
			},
			expectedResult: entity.Album{ID: entity.AlbumID("1"), Title: "Test Album"},
			expectError:    false,
		},
		{
			name: "GetAlbumByID_NotFound",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, title FROM album WHERE id = ?")).
					WithArgs("1").
					WillReturnError(sql.ErrNoRows)
			},
			action: func(r *mysql.AlbumRepository) interface{} {
				_, err := r.GetAlbumByID(context.Background(), "1")
				return err
			},
			expectedResult: nil,
			expectError:    true,
			expectedErr:    customerr.ErrAlbumNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock DB
			db, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo, err := mysql.NewAlbumRepository(db)
			assert.NoError(t, err)

			// Execute action
			result := tt.action(repo)

			// Assertions
			if tt.expectError {
				assert.Error(t, result.(error))
				assert.EqualError(t, result.(error), tt.expectedErr.Error())
			} else {
				assert.Equal(t, tt.expectedResult, result)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Add missing helper functions that were assumed in the original code
func BuildAlbumEntity(dbAlbum mysql.Album) entity.Album {
	return entity.Album{
		ID:    entity.AlbumID(dbAlbum.ID),
		Title: dbAlbum.Title,
	}
}

func BuildDBAlbum(entity entity.Album) mysql.Album {
	return mysql.Album{
		ID:    string(entity.ID),
		Title: entity.Title,
	}
}
