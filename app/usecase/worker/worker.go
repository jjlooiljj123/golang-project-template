package worker

import (
	"boilerplate/app/infrastructure/sqs/queue"
	"context"
	"errors"
	"log"
	"sync/atomic"
	"time"
)

// Define a custom error for panic recovery
var (
	errPanicDefaultMessage = errors.New("panic error from processor")
)

// MessageHandler defines the type for functions that handle messages
type MessageHandler func(context.Context, *queue.Queue)

// Processor defines the interface for processing logic
type Processor interface {
	Process(ctx context.Context, done chan<- struct{}, handler MessageHandler)
}

// Worker manages multiple instances of a Processor
type Worker struct {
	processor Processor

	// Number of goroutines to create for processor
	goroutinesNumber int

	// Duration to wait before retry when a processor goroutine panics
	retryInterval time.Duration

	// Duration to wait before force quitting after receiving stop signal from upstream
	waitTime time.Duration

	// Channel to send/receive error when a processor goroutine panics
	panicErrCh chan error

	// Channel to signal upstream that worker is done
	done chan struct{}

	// Channel for processor goroutines to signal worker that they have done their job
	processorDone chan struct{}
}

// NewWorker returns a new instance of Worker
func NewWorker(processor Processor, goroutinesNumber int, retryInterval time.Duration, waitTime time.Duration) *Worker {
	return &Worker{
		processor:        processor,
		goroutinesNumber: goroutinesNumber,
		retryInterval:    retryInterval,
		waitTime:         waitTime,
		panicErrCh:       make(chan error),
		done:             make(chan struct{}),
		processorDone:    make(chan struct{}, goroutinesNumber),
	}
}

// Start initiates the worker with the specified number of goroutines
func (w *Worker) Start(ctx context.Context, handler MessageHandler) {
	var processorDoneNumber int32

	// Goroutine to handle retries on panic
	go func() {
		for {
			select {
			case <-w.panicErrCh:
				time.Sleep(w.retryInterval)
				w.startProcess(ctx, handler)
			case <-w.done:
				return
			}
		}
	}()

	// Goroutine to manage graceful shutdown
	go func() {
		<-ctx.Done()
		for {
			select {
			case <-w.processorDone:
				atomic.AddInt32(&processorDoneNumber, 1)
				if processorDoneNumber == int32(w.goroutinesNumber) {
					close(w.done)
					return
				}
			case <-time.After(w.waitTime):
				close(w.done)
				return
			}
		}
	}()

	// Start the specified number of processing goroutines
	for i := 0; i < w.goroutinesNumber; i++ {
		w.startProcess(ctx, handler)
	}
}

// Done returns the worker's done channel for upstream to listen to
func (w *Worker) Done() <-chan struct{} {
	return w.done
}

// startProcess starts a single processor goroutine with recovery logic
func (w *Worker) startProcess(ctx context.Context, handler MessageHandler) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					// err = errors.Wrap(r, errPanicDefaultMessage.Error())
					err = r
				default:
					err = errPanicDefaultMessage
				}

				log.Printf("worker panic")
				w.panicErrCh <- err
			}
		}()

		w.processor.Process(ctx, w.processorDone, handler)
	}()
}
