LOCALSTACK_IP="http://localhost:4566"
ALBUM_QUEUE="http://sqs.us-east-1.localhost.localstack.cloud:4566/000000000000/album"

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    echo "Please enter the SQS json file"
    exit 1
fi

aws --endpoint-url $LOCALSTACK_IP sqs send-message \
  --queue-url $ALBUM_QUEUE \
  --message-body file://$1




