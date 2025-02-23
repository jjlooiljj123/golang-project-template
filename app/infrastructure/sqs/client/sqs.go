package sqs

import (
	"context"
	"net/http"

	infraConfig "boilerplate/app/infrastructure/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// NewSQSClient initializes and returns an SQS client with the given configuration
func NewSQSClient(ctx context.Context, workerConfig infraConfig.WorkerConfig) (*sqs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(workerConfig.Region),
		config.WithHTTPClient(&http.Client{
			Timeout: workerConfig.SQS.HTTPTimeout, // Example timeout, adjust as needed
		}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			workerConfig.AccessKeyID,
			workerConfig.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	// Set the base endpoint after loading the config
	cfg.BaseEndpoint = aws.String(workerConfig.SQSHost)

	return sqs.NewFromConfig(cfg), nil
}
