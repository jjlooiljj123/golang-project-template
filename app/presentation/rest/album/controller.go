package album

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"boilerplate/app/domain/dto"
	"boilerplate/app/domain/entity"
	"boilerplate/app/domain/errors"

	albumservice "boilerplate/app/usecase/interface"
)

type Controller struct {
	albumService albumservice.AlbumInterface
}

func NewController(
	albumService albumservice.AlbumInterface,
) *Controller {
	return &Controller{
		albumService: albumService,
	}
}

func (c *Controller) GetAlbumsHandler(ctx *gin.Context) {
	// headers, ok := middleware.GetCommonHeadersFromContext(ctx.Request.Context())
	// if ok {
	// 	log.Printf("Request ID: %s, User-Agent: %s", headers.RequestID, headers.UserAgent)
	// }
	data, err := c.albumService.GetAllAlbums(ctx)
	if err != nil {
		// Handle the error
		c.handleError(ctx, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, data)
}

// CreateAlbumHandler handles POST requests to create a album
func (c *Controller) CreateAlbumHandler(ctx *gin.Context) {
	var album dto.Album
	if err := ctx.ShouldBindJSON(&album); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	entityAlbum := entity.Album{
		ID:    entity.AlbumID(album.ID),
		Title: album.Title,
	}

	id, err := c.albumService.CreateAlbum(ctx, entityAlbum)
	if err != nil {
		c.handleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id, "message": "Album created successfully"})
}

func (c *Controller) GetAlbumByIDHandler(ctx *gin.Context) {
	id := ctx.Param("id")
	album, err := c.albumService.GetAlbumByID(ctx, id)
	if err != nil {
		if err.Error() == "album not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		} else {
			c.handleError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, album)
}

func (c *Controller) GetJsonPostHandler(ctx *gin.Context) {
	data, err := c.albumService.GetFromThirdPartyAPI(ctx)
	if err != nil {
		// Handle the error
		c.handleError(ctx, err)
		return
	}
	ctx.IndentedJSON(http.StatusOK, data)
}

// handleError is a helper method to manage error responses
func (c *Controller) handleError(ctx *gin.Context, err error) {
	// Determine the HTTP status code based on the error type
	status := http.StatusInternalServerError
	message := "Internal Server Error"

	if errors.IsAlbumNotFound(err) {
		status = http.StatusNotFound
		message = "Album not found"
	} else if errors.IsInvalidInput(err) {
		status = http.StatusBadRequest
		message = "Invalid input"
	}

	// Log the error for debugging purposes
	ctx.Error(err)

	// Respond with the appropriate status and message
	ctx.JSON(status, gin.H{"error": message})
}
