package queue

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Queue represents an SQS queue
type Queue struct {
	client   *sqs.Client
	queueURL string
}

// NewQueue creates a new Queue instance
func NewQueue(client *sqs.Client, queueURL string) *Queue {
	return &Queue{
		client:   client,
		queueURL: queueURL,
	}
}

// SendMessage sends a single message to the queue
func (q *Queue) SendMessage(message string) error {
	_, err := q.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(q.queueURL),
		MessageBody: aws.String(message),
	})
	return err
}

// SendMessages sends multiple messages to the queue in batch
func (q *Queue) SendMessages(messages []string) error {
	if len(messages) == 0 {
		return nil
	}

	entries := make([]types.SendMessageBatchRequestEntry, len(messages))
	for i, m := range messages {
		entries[i] = types.SendMessageBatchRequestEntry{
			Id:          aws.String(strconv.Itoa(i)), // Simple way to ensure uniqueness
			MessageBody: aws.String(m),
		}
	}
	_, err := q.client.SendMessageBatch(context.TODO(), &sqs.SendMessageBatchInput{
		QueueUrl: aws.String(q.queueURL),
		Entries:  entries,
	})
	return err
}

// MaxNumberOfSqsMessageForRead is the maximum number of messages that can be read from SQS in one request
const MaxNumberOfSqsMessageForRead = 10

// ReceiveMessages retrieves messages from the queue
func (q *Queue) ReceiveMessages(numberOfMessages int, waitTime time.Duration) ([]types.Message, error) {
	if numberOfMessages > MaxNumberOfSqsMessageForRead {
		numberOfMessages = MaxNumberOfSqsMessageForRead
	}
	output, err := q.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueURL),
		MaxNumberOfMessages: int32(numberOfMessages),
		WaitTimeSeconds:     int32(waitTime.Seconds()),
	})
	if err != nil {
		return nil, err
	}
	return output.Messages, nil
}

// DeleteMessage deletes a message from the queue
func (q *Queue) DeleteMessage(receiptHandle string) error {
	_, err := q.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

// GetNumberOfQueueMessages retrieves the approximate number of messages in the queue
func (q *Queue) GetNumberOfQueueMessages() (int, error) {
	output, err := q.client.GetQueueAttributes(context.TODO(), &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(q.queueURL),
		AttributeNames: []types.QueueAttributeName{types.QueueAttributeNameApproximateNumberOfMessages},
	})
	if err != nil {
		return 0, err
	}
	if val, ok := output.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessages)]; ok {
		return strconv.Atoi(val)
	}
	return 0, nil
}

// ChangeMessageVisibility changes the visibility timeout for a message
func (q *Queue) ChangeMessageVisibility(receiptHandle string, newTimeout int32) error {
	_, err := q.client.ChangeMessageVisibility(context.TODO(), &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(q.queueURL),
		ReceiptHandle:     aws.String(receiptHandle),
		VisibilityTimeout: newTimeout,
	})
	return err
}
