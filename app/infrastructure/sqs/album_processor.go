package sqs

import (
	"context"
	"log"
	"time"

	infraConfig "boilerplate/app/infrastructure/config"
	services "boilerplate/app/usecase/interface"

	"github.com/aws/aws-sdk-go-v2/service/sqs"

	"boilerplate/app/infrastructure/sqs/queue"
	"boilerplate/app/usecase/worker"
)

// AlbumProcessor is an implementation of the Processor interface for SQS
type AlbumProcessor struct {
	client   *sqs.Client
	config   infraConfig.WorkerConfig
	queueURL string

	service services.MessageProcessor
}

// NewAlbumProcessor creates a new AlbumProcessor
func NewAlbumProcessor(client *sqs.Client, queueURL string, workerConfig infraConfig.WorkerConfig, svc services.MessageProcessor) *AlbumProcessor {
	return &AlbumProcessor{
		client:   client,
		queueURL: queueURL,
		config:   workerConfig,
		service:  svc,
	}
}

// Process implements the Processor interface for handling SQS messages
func (p *AlbumProcessor) Process(ctx context.Context, done chan<- struct{}, handler worker.MessageHandler) {
	q := queue.NewQueue(p.client, p.queueURL)

	for {
		select {
		case <-ctx.Done():
			done <- struct{}{}
			return
		default:
			handler(ctx, q)
		}
	}
}

// DefaultMessageHandler provides the default behavior for processing messages
func (p *AlbumProcessor) DefaultMessageHandler(ctx context.Context, q *queue.Queue) {
	messages, err := q.ReceiveMessages(queue.MaxNumberOfSqsMessageForRead, p.config.SQS.LongPollingWait)
	if err != nil {
		// Handle error, perhaps with logging
		return
	}

	for _, message := range messages {
		// Process message using the injected service
		if err := p.service.ProcessMessage(ctx, *message.Body); err != nil {
			// Handle processing error
			continue
		}

		// Delete the message after processing
		if err := q.DeleteMessage(*message.ReceiptHandle); err != nil {
			log.Printf("Failed to delete message %s: %v", *message.MessageId, err)
		}
	}

	// If no messages were received, wait before polling again
	if len(messages) == 0 {
		time.Sleep(1 * time.Second)
	}
}
