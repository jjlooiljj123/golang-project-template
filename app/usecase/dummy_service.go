package services

import (
	"context"
	"fmt"
)

// DummyService represents a simple service for demonstration purposes
type DummyService struct {
}

// ProcessMessage simulates processing a message from SQS
func (s *DummyService) ProcessMessage(ctx context.Context, message string) error {
	fmt.Printf("DummyService processing message: %s\n", message)
	return nil
}

// NewDummyService creates a new instance of DummyService
func NewDummyService() *DummyService {
	return &DummyService{}
}
