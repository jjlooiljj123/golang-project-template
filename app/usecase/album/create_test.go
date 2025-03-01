package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"boilerplate/app/domain/entity"
	"boilerplate/app/infrastructure/repositories/interface/mocks"
	albumservice "boilerplate/app/usecase/album"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAlbumRepository is a mock implementation of the album repository
type MockAlbumRepository struct {
	mock.Mock
}

// CreateAlbum mock implementation
func (m *MockAlbumRepository) CreateAlbum(ctx context.Context, album entity.Album) (string, error) {
	args := m.Called(ctx, album)
	return args.String(0), args.Error(1)
}

func TestService_CreateAlbum(t *testing.T) {
	// Test cases
	tests := []struct {
		name          string
		album         entity.Album
		mockID        string
		mockError     error
		expectedID    string
		expectedError error
	}{
		{
			name: "Successful creation",
			album: entity.Album{
				ID:    entity.AlbumID("album-123"),
				Title: "Test Album",
			},
			mockID:        "album-123",
			mockError:     nil,
			expectedID:    "album-123",
			expectedError: nil,
		},
		{
			name: "Repository error",
			album: entity.Album{
				ID:    entity.AlbumID("album-456"),
				Title: "Test Album",
			},
			mockID:        "",
			mockError:     errors.New("database error"),
			expectedID:    "",
			expectedError: errors.New("service error creating album: database error"),
		},
		{
			name: "Empty album",
			album: entity.Album{
				ID:    entity.AlbumID(""),
				Title: "",
			},
			mockID:        "album-empty",
			mockError:     nil,
			expectedID:    "album-empty",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := mocks.NewRepositoryInterface(t)

			// Set up mock expectations
			mockRepo.On("CreateAlbum", mock.Anything, tt.album).
				Return(tt.mockID, tt.mockError).
				Once()

			// Create service with mock repository
			service := albumservice.NewService(
				mockRepo,
				nil,
				0*time.Second,
				nil,
			)

			// Call the method
			ctx := context.Background()
			id, err := service.CreateAlbum(ctx, tt.album)

			// Assertions
			assert.Equal(t, tt.expectedID, id, "ID should match expected")

			if tt.expectedError != nil {
				assert.Error(t, err, "Should return an error")
				assert.EqualError(t, err, tt.expectedError.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Should not return an error")
			}

			// Verify all expectations were met
			mockRepo.AssertExpectations(t)
		})
	}
}
