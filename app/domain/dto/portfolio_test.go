package dto_test

import (
	"testing"

	"github.com/google/uuid"

	"boilerplate/app/domain/dto"
	"boilerplate/app/domain/entity"
)

func TestBuildAlbumDTO(t *testing.T) {
	tests := []struct {
		name     string
		entity   entity.Album
		expected dto.Album
	}{
		{
			name: "Typical case",
			entity: entity.Album{
				ID:    entity.AlbumID(uuid.New().String()),
				Title: "Test Album",
			},
			expected: dto.Album{
				ID:    "", // Will be set in test
				Title: "Test Album",
			},
		},
		{
			name: "Empty title",
			entity: entity.Album{
				ID:    entity.AlbumID(uuid.New().String()),
				Title: "",
			},
			expected: dto.Album{
				ID:    "", // Will be set in test
				Title: "",
			},
		},
		{
			name: "Long title",
			entity: entity.Album{
				ID:    entity.AlbumID(uuid.New().String()),
				Title: "This is a very long album title that should still work",
			},
			expected: dto.Album{
				ID:    "", // Will be set in test
				Title: "This is a very long album title that should still work",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the expected ID based on the entity ID
			tt.expected.ID = tt.entity.ID.String()

			result := dto.BuildAlbumDTO(tt.entity)

			if result.ID != tt.expected.ID {
				t.Errorf("Test %s: Expected ID %s, got %s", tt.name, tt.expected.ID, result.ID)
			}

			if result.Title != tt.expected.Title {
				t.Errorf("Test %s: Expected Title %s, got %s", tt.name, tt.expected.Title, result.Title)
			}
		})
	}
}
