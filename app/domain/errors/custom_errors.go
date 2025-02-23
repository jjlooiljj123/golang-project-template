package errors

import "errors"

// Custom error definitions for specific scenarios
var (
	ErrAlbumNotFound  = errors.New("album not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrInternalServer = errors.New("internal server error")
	// Add more custom errors here as needed
)

// IsAlbumNotFound checks if the error is a album not found error
func IsAlbumNotFound(err error) bool {
	return err == ErrAlbumNotFound
}

// IsInvalidInput checks if the error is an invalid input error
func IsInvalidInput(err error) bool {
	return err == ErrInvalidInput
}

// IsInternalServer checks if the error is an internal server error
func IsInternalServer(err error) bool {
	return err == ErrInternalServer
}
