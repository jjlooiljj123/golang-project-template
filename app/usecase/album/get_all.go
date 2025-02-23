package services

import (
	"boilerplate/app/domain/dto"
	"context"
)

func (s *Service) GetAllAlbums(ctx context.Context) ([]dto.Album, error) {
	albums, err := s.albumRepo.GetAlbums(ctx)
	if err != nil {
		return []dto.Album{}, err
	}

	albumDTOs := make([]dto.Album, len(albums))
	for i, albumEntity := range albums {
		dto := dto.BuildAlbumDTO(albumEntity)
		albumDTOs[i] = dto
	}

	return albumDTOs, nil
}
