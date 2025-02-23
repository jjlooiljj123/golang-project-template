package mysql

import (
	"boilerplate/app/domain/entity"
	"context"
)

// RepositoryInterface defines the interface for MySQL operations
type RepositoryInterface interface {
	GetAlbums(ctx context.Context) ([]entity.Album, error)
	CreateAlbum(ctx context.Context, album entity.Album) (string, error)
	GetAlbumByID(ctx context.Context, id string) (entity.Album, error)
}
