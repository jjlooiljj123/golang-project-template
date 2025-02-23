package service

import (
	"boilerplate/app/domain/dto"
	"boilerplate/app/domain/entity"
	"context"
)

type GetAlbumInterface interface {
	GetAllAlbums(ctx context.Context) ([]dto.Album, error)
	GetAlbumByID(ctx context.Context, id string) (dto.Album, error)
}

type CreateAlbumInterface interface {
	CreateAlbum(ctx context.Context, album entity.Album) (string, error)
}

type GetJSONPostInterface interface {
	GetFromThirdPartyAPI(ctx context.Context) ([]dto.Post, error)
}

type AlbumInterface interface {
	GetAlbumInterface
	CreateAlbumInterface
	GetJSONPostInterface
}
