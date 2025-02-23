package services

import (
	"time"

	httpClientJsonPostInterface "boilerplate/app/infrastructure/httpclient/interface"
	"boilerplate/app/infrastructure/redis"
	albumsRepositories "boilerplate/app/infrastructure/repositories/interface"
)

type Service struct {
	albumRepo albumsRepositories.RepositoryInterface
	cache     *redis.RedisCache
	cacheT    time.Duration // Duration for cache expiration

	jsonPostService httpClientJsonPostInterface.HttpClientJsonPostInterface
}

func NewService(
	albumRepo albumsRepositories.RepositoryInterface,
	redisCache *redis.RedisCache,
	cacheExpiration time.Duration,
	jsonPostService httpClientJsonPostInterface.HttpClientJsonPostInterface,
) *Service {
	return &Service{
		albumRepo:       albumRepo,
		cache:           redisCache,
		cacheT:          cacheExpiration,
		jsonPostService: jsonPostService,
	}
}
