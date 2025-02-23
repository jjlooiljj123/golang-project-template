package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"boilerplate/app/domain/entity"
	"boilerplate/app/domain/errors"
)

// Album represents the structure of a album in our application with db tags for column mapping
type Album struct {
	ID    string `db:"id"`
	Title string `db:"title"`
}

// AlbumRepository implements RepositoryInterface for MySQL database operations
type AlbumRepository struct {
	db *sql.DB
}

// NewAlbumRepository initializes a new MySQL repository with the given connection string
func NewAlbumRepository(db *sql.DB) (*AlbumRepository, error) {
	return &AlbumRepository{db: db}, nil
}

func (r *AlbumRepository) GetAlbums(ctx context.Context) ([]entity.Album, error) {
	var albums []entity.Album
	var dbAlbums []Album
	rows, err := r.db.QueryContext(ctx, "SELECT id, title FROM album")
	if err != nil {
		return nil, fmt.Errorf("error querying data: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		dbAlbums = append(dbAlbums, album)
	}

	// Check for errors from iterating over rows.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %v", err)
	}

	for _, dbPordbAlbum := range dbAlbums {
		album := BuildAlbumEntity(dbPordbAlbum)
		albums = append(albums, album)
	}

	return albums, nil
}

// CreateAlbum inserts a new album into the database
func (r *AlbumRepository) CreateAlbum(ctx context.Context, entity entity.Album) (string, error) {

	album := BuildDBAlbum(entity)

	// Insert the new album into the database
	stmt, err := r.db.Prepare("INSERT INTO album (id, title) VALUES (?, ?)")
	if err != nil {
		return "", fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, album.ID, album.Title)
	if err != nil {
		return "", fmt.Errorf("error executing insert: %v", err)
	}

	return album.ID, nil
}

// GetAlbumByID retrieves a specific album by its ID from MySQL
func (r *AlbumRepository) GetAlbumByID(ctx context.Context, id string) (entity.Album, error) {
	var album Album
	err := r.db.QueryRowContext(ctx, "SELECT id, title FROM album WHERE id = ?", id).Scan(&album.ID, &album.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Album{}, errors.ErrAlbumNotFound
		}
		return entity.Album{}, errors.ErrInternalServer
	}
	entityAlbum := BuildAlbumEntity(album)
	return entityAlbum, nil
}
