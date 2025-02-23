package services

import (
	"context"
	"fmt"

	httpclient "boilerplate/app/infrastructure/httpclient/interface"
)

// DummyService represents a simple service for demonstration purposes
type DummyService struct {
	jsonPostClient httpclient.HttpClientJsonPostInterface
}

// ProcessMessage simulates processing a message from SQS
func (s *DummyService) ProcessMessage(ctx context.Context, message string) error {
	fmt.Printf("DummyService processing message: %s\n", message)
	posts, err := s.jsonPostClient.GetPosts(ctx)
	if err != nil {
		fmt.Printf("err getting http posts: %v", err)
	}
	fmt.Printf("done getting http posts: %v", posts)
	return nil
}

// NewDummyService creates a new instance of DummyService
func NewDummyService() *DummyService {
	return &DummyService{}
}
