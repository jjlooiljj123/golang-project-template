#!/bin/sh

echo "LocalStack SQS is up, creating 'album' queue..."

# Create the 'album' queue
awslocal sqs create-queue --queue-name album

echo "Queue 'album' created."