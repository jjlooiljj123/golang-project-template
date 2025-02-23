package httpclient

import (
	"boilerplate/app/domain/entity"
	"context"
)

type HttpClientJsonPostInterface interface {
	GetPosts(ctx context.Context) ([]entity.Post, error)
}
