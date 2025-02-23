package mysql

import "boilerplate/app/domain/entity"

func BuildAlbumEntity(album Album) entity.Album {
	return entity.Album{
		ID:    entity.AlbumID(album.ID),
		Title: album.Title,
	}
}
