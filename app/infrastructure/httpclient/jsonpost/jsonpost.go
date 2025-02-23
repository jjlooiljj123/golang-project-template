package jsonpost

import (
	"boilerplate/app/domain/entity"
	"boilerplate/app/infrastructure/config"
	"boilerplate/app/infrastructure/httpclient"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HttpJsonPost struct {
	http      *httpclient.Client
	appConfig *config.AppConfig
}

func NewHttpJsonPost(httpClient *httpclient.Client, appConfig *config.AppConfig) *HttpJsonPost {
	return &HttpJsonPost{
		http:      httpClient,
		appConfig: appConfig,
	}
}

type JsonPost struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (s *HttpJsonPost) GetPosts(ctx context.Context) ([]entity.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, s.appConfig.APITimeout)
	defer cancel()

	url := s.appConfig.JSONPlaceHolderURL + "/posts"
	resp, err := s.http.Get(ctx, url, nil)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, errors.New("request timeout")
		}
		return nil, fmt.Errorf("error fetching posts: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := httpclient.ReadBody(resp)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var posts []JsonPost
	if err := json.Unmarshal(body, &posts); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	results := make([]entity.Post, len(posts))
	for i, post := range posts {
		results[i] = entity.Post{
			UserID: post.UserID,
			ID:     post.ID,
			Title:  post.Title,
			Body:   post.Body,
		}
	}

	return results, nil
}
