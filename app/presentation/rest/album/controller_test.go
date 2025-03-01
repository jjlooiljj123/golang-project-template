package album_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"boilerplate/app/domain/dto"
	"boilerplate/app/domain/entity"
	customerr "boilerplate/app/domain/errors"
	"boilerplate/app/presentation/rest/album"
	"boilerplate/app/usecase/interface/mocks"
)

func TestController(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMock      func(*mocks.AlbumInterface)
		method         string
		url            string
		body           interface{}
		expectedStatus int
		expectedBody   string // Changed to string
	}{
		// GetAlbumsHandler tests
		{
			name: "GetAlbums_Success",
			setupMock: func(m *mocks.AlbumInterface) {
				m.On("GetAllAlbums", mock.Anything).Return([]dto.Album{{ID: "1", Title: "Test"}}, nil)
			},
			method:         "GET",
			url:            "/albums",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":"1","title":"Test"}]`,
		},
		{
			name: "GetAlbums_Error",
			setupMock: func(m *mocks.AlbumInterface) {
				m.On("GetAllAlbums", mock.Anything).Return([]dto.Album{}, errors.New("db error"))
			},
			method:         "GET",
			url:            "/albums",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal Server Error"}`,
		},

		// CreateAlbumHandler tests
		{
			name: "CreateAlbum_Success",
			setupMock: func(m *mocks.AlbumInterface) {
				m.On("CreateAlbum", mock.Anything, entity.Album{ID: entity.AlbumID("1"), Title: "Test"}).
					Return("1", nil)
			},
			method:         "POST",
			url:            "/albums",
			body:           dto.Album{ID: "1", Title: "Test"},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"id":"1","message":"Album created successfully"}`,
		},
		{
			name:           "CreateAlbum_InvalidInput",
			method:         "POST",
			url:            "/albums",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid input"}`,
		},

		// GetAlbumByIDHandler tests
		{
			name: "GetAlbumByID_Success",
			setupMock: func(m *mocks.AlbumInterface) {
				m.On("GetAlbumByID", mock.Anything, "1").Return(dto.Album{ID: "1", Title: "Test"}, nil)
			},
			method:         "GET",
			url:            "/albums/1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":"1","title":"Test"}`,
		},
		{
			name: "GetAlbumByID_NotFound",
			setupMock: func(m *mocks.AlbumInterface) {
				m.On("GetAlbumByID", mock.Anything, "1").Return(dto.Album{}, customerr.ErrAlbumNotFound)
			},
			method:         "GET",
			url:            "/albums/1",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Album not found"}`,
		},

		// GetJsonPostHandler tests
		{
			name: "GetJsonPost_Success",
			setupMock: func(m *mocks.AlbumInterface) {
				m.On("GetFromThirdPartyAPI", mock.Anything).Return([]dto.Post{
					{UserID: 1, ID: 1, Title: "title_1", Body: "body_1"},
				}, nil)
			},
			method:         "GET",
			url:            "/json",
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"body":"body_1", "id":1, "title":"title_1", "userId":1}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := mocks.NewAlbumInterface(t)
			if tt.setupMock != nil {
				tt.setupMock(mockService)
			}

			// Create controller
			controller := album.NewController(mockService)

			// Set up gin router
			router := gin.New()
			router.GET("/albums", controller.GetAlbumsHandler)
			router.POST("/albums", controller.CreateAlbumHandler)
			router.GET("/albums/:id", controller.GetAlbumByIDHandler)
			router.GET("/json", controller.GetJsonPostHandler)

			// Create request
			var bodyBytes []byte
			if tt.body != nil {
				if str, ok := tt.body.(string); ok {
					bodyBytes = []byte(str)
				} else {
					bodyBytes, _ = json.Marshal(tt.body)
				}
			}
			req, _ := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(bodyBytes))
			if tt.method == "POST" {
				req.Header.Set("Content-Type", "application/json")
			}

			// Create response recorder
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Compare response body as string
			// Remove whitespace and newlines for consistent comparison
			actualBody := w.Body.String()
			assert.JSONEq(t, tt.expectedBody, actualBody, "Response body should match expected JSON")

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
