package services

import (
	"boilerplate/app/domain/dto"
	"context"
)

// GetAlbumByID retrieves a album by ID using the repository
func (s *Service) GetFromThirdPartyAPI(ctx context.Context) ([]dto.Post, error) {
	entityPosts, err := s.jsonPostService.GetPosts(ctx)
	if err != nil {
		return []dto.Post{}, err
	}
	results := make([]dto.Post, len(entityPosts))
	for i, post := range entityPosts {
		results[i] = dto.Post{
			UserID: post.UserID,
			ID:     post.ID,
			Title:  post.Title,
			Body:   post.Body,
		}
	}
	return results, nil
}
