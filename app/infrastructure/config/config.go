package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	MySQLHost     string `env:"MYSQL_HOST"`
	MySQLUser     string `env:"MYSQL_USER"`
	MySQLPassword string `env:"MYSQL_PASSWORD"`
	MySQLDatabase string `env:"MYSQL_DATABASE"`

	RedisHost     string        `env:"REDIS_HOST"`
	RedisPort     string        `env:"REDIS_PORT"`
	CacheDuration time.Duration `env:"CACHE_DURATION"`

	JSONPlaceHolderURL string        `env:"JSON_PLACEHOLDER_URL"`
	APITimeout         time.Duration `env:"API_TIMEOUT"`
	HandlerTimeout     time.Duration `env:"HANDLER_TIMEOUT"`
}

var AppCfg AppConfig

// LoadConfig reads environment variables from a .env file into the appConfig struct
func LoadConfig() error {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Map environment variables to appConfig struct fields
	AppCfg.MySQLHost = os.Getenv("MYSQL_HOST")
	AppCfg.MySQLUser = os.Getenv("MYSQL_USER")
	AppCfg.MySQLPassword = os.Getenv("MYSQL_PASSWORD")
	AppCfg.MySQLDatabase = os.Getenv("MYSQL_DATABASE")
	AppCfg.RedisHost = os.Getenv("REDIS_HOST")
	AppCfg.RedisPort = os.Getenv("REDIS_PORT")

	// Parse duration for cache, default to 5 minutes if not set or invalid
	cacheDurationStr := os.Getenv("CACHE_DURATION")
	if cacheDurationStr == "" {
		AppCfg.CacheDuration = 5 * time.Minute
	} else {
		duration, err := time.ParseDuration(cacheDurationStr)
		if err != nil {
			return fmt.Errorf("invalid CACHE_DURATION format: %v", err)
		}
		AppCfg.CacheDuration = duration
	}

	AppCfg.JSONPlaceHolderURL = os.Getenv("JSON_PLACEHOLDER_URL")
	apiTimeoutStr := os.Getenv("API_TIMEOUT")
	if apiTimeoutStr == "" {
		AppCfg.APITimeout = 5 * time.Second
	} else {
		duration, err := time.ParseDuration(apiTimeoutStr)
		if err != nil {
			return fmt.Errorf("invalid API_TIMEOUT format: %v", err)
		}
		AppCfg.APITimeout = duration
	}
	handlerTimeoutStr := os.Getenv("HANDLER_TIMEOUT")
	if handlerTimeoutStr == "" {
		AppCfg.HandlerTimeout = 15 * time.Second
	} else {
		duration, err := time.ParseDuration(handlerTimeoutStr)
		if err != nil {
			return fmt.Errorf("invalid HANDLER_TIMEOUT format: %v", err)
		}
		AppCfg.HandlerTimeout = duration
	}

	return nil
}
