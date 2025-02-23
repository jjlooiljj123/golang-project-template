package dto

import "boilerplate/app/domain/entity"

type Album struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func BuildAlbumDTO(albumEntity entity.Album) Album {
	return Album{
		ID:    albumEntity.ID.String(),
		Title: albumEntity.Title,
	}
}
