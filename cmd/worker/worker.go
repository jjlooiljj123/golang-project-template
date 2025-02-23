package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	infraConfig "boilerplate/app/infrastructure/config"
	"boilerplate/app/infrastructure/sqs"
	sqsclient "boilerplate/app/infrastructure/sqs/client"
	"boilerplate/app/infrastructure/sqs/queue"
	services "boilerplate/app/usecase"
	"boilerplate/app/usecase/worker"
)

func main() {
	// Load worker configuration
	workerConfig, err := infraConfig.LoadWorkerConfig()
	if err != nil {
		log.Fatalf("Failed to load worker configuration: %v", err)
	}

	// Initialize SQS client
	sqsClient, err := sqsclient.NewSQSClient(context.TODO(), workerConfig)
	if err != nil {
		log.Fatalf("Failed to initialize SQS client: %v", err)
	}

	// Create the dummy service
	dummyService := services.NewDummyService()

	// Setup the processor and worker
	sqsProcessor := sqs.NewAlbumProcessor(sqsClient, workerConfig.QueueURL, workerConfig, dummyService)

	// Define the handler, potentially wrapped with middleware
	handler := func(ctx context.Context, q *queue.Queue) {
		sqsProcessor.DefaultMessageHandler(ctx, q)
	}

	// Custom middleware function
	middleware := func(next worker.MessageHandler) worker.MessageHandler {
		return func(ctx context.Context, q *queue.Queue) {
			// Middleware logic before message processing
			log.Println("Middleware: Before message processing")
			next(ctx, q)
			// Middleware logic after message processing
			log.Println("Middleware: After message processing")
		}
	}

	// Wrap the default handler with middleware
	wrappedHandler := middleware(handler)

	// Pass the wrapped handler to Process
	albumWorker := worker.NewWorker(
		sqsProcessor,
		workerConfig.AlbumWorker.GoroutinesNumber,
		workerConfig.AlbumWorker.RetryInterval,
		workerConfig.AlbumWorker.WaitTime,
	)

	var wg sync.WaitGroup
	// Create a context for the worker
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the worker
	wg.Add(1)
	albumWorker.Start(ctx, wrappedHandler)

	// Wait for termination signal for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Signal the worker to stop
	cancel()

	// Wait for worker to finish
	go func(albumWorker *worker.Worker) {
		<-albumWorker.Done()
		wg.Done()
	}(albumWorker)

	wg.Wait()
	log.Println("Worker has stopped.")
}
