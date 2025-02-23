package services

import (
	"boilerplate/app/domain/dto"
	"boilerplate/app/domain/entity"
	"boilerplate/app/domain/errors"
	"context"
	"encoding/json"
	"fmt"
)

// GetAlbumByID retrieves a album by ID using the repository
func (s *Service) GetAlbumByID(ctx context.Context, id string) (dto.Album, error) {

	var album entity.Album
	// Try to get from cache
	if s.cache != nil {
		cachedData, err := s.cache.GetFromCache(id)
		if err == nil {
			if jsonErr := json.Unmarshal(cachedData, &album); jsonErr == nil {
				return dto.BuildAlbumDTO(album), nil
			}
		}
	}

	album, err := s.albumRepo.GetAlbumByID(ctx, id)
	if err != nil {
		if errors.IsAlbumNotFound(err) {
			return dto.Album{}, errors.ErrAlbumNotFound
		}

		return dto.Album{}, fmt.Errorf("service error getting album: %v", err)
	}

	// Store in cache for next time
	if s.cache != nil {
		if err := s.cache.SetToCache(id, album, s.cacheT); err != nil {
			// Log cache error but don't fail the request if cache write fails
			fmt.Printf("Failed to cache album %s: %v\n", id, err)
		}
	}

	dto := dto.BuildAlbumDTO(album)
	return dto, nil
}
