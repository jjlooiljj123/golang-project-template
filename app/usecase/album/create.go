package services

import (
	"boilerplate/app/domain/entity"
	"context"
	"fmt"
)

// CreateAlbum creates a new album using the repository
func (s *Service) CreateAlbum(ctx context.Context, album entity.Album) (string, error) {
	id, err := s.albumRepo.CreateAlbum(ctx, album)
	if err != nil {
		return "", fmt.Errorf("service error creating album: %v", err)
	}
	return id, nil
}
