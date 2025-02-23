package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// WorkerConfig holds configuration for the worker, including AWS SQS settings
type WorkerConfig struct {
	QueueURL        string
	Region          string
	SQSHost         string
	AccessKeyID     string
	SecretAccessKey string

	SQS         SQSConfig
	AlbumWorker AlbumWorkerConfig
}

// AlbumWorkerConfig holds configuration specific to the album worker
type AlbumWorkerConfig struct {
	GoroutinesNumber int
	RetryInterval    time.Duration
	WaitTime         time.Duration
}

// SQSConfig holds SQS-specific configurations
type SQSConfig struct {
	HTTPTimeout     time.Duration
	LongPollingWait time.Duration
}

// LoadWorkerConfig loads worker configuration from a .env file
func LoadWorkerConfig() (WorkerConfig, error) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		return WorkerConfig{}, fmt.Errorf("error loading .env file: %v", err)
	}

	config := WorkerConfig{
		QueueURL:        os.Getenv("SQS_QUEUE_URL"),
		Region:          os.Getenv("AWS_REGION"),
		SQSHost:         os.Getenv("AWS_SQS_HOST"),
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SQS: SQSConfig{
			HTTPTimeout:     parseDurationEnv("SQS_HTTP_TIMEOUT", 30*time.Second),
			LongPollingWait: parseDurationEnv("SQS_LONG_POLLING_WAIT", 20*time.Second),
		},
		AlbumWorker: AlbumWorkerConfig{
			GoroutinesNumber: parseIntEnv("ALBUM_WORKER_GOROUTINES", 1),
			RetryInterval:    parseDurationEnv("ALBUM_WORKER_RETRY_INTERVAL", 5*time.Second),
			WaitTime:         parseDurationEnv("ALBUM_WORKER_WAIT_TIME", 10*time.Second),
		},
	}

	if config.QueueURL == "" ||
		config.Region == "" ||
		config.SQSHost == "" ||
		config.AccessKeyID == "" ||
		config.SecretAccessKey == "" {
		return WorkerConfig{}, fmt.Errorf("missing required environment variables in .env file")
	}

	return config, nil
}

// Helper functions to parse environment variables
func parseIntEnv(key string, defaultValue int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}

func parseDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}
