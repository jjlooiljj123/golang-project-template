package mysql

import "boilerplate/app/domain/entity"

func BuildDBAlbum(entity entity.Album) Album {
	return Album{
		ID:    entity.ID.String(),
		Title: entity.Title,
	}
}
