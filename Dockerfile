# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.23-alpine as builder

# Create and change to the app directory.
WORKDIR /app

# Copy go mod and sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container
COPY . .

# Build the Go app
RUN go build -o main ./cmd/album

# Build the worker
RUN go build -o worker ./cmd/worker

# Use a minimal image to run the binary
FROM alpine:3.14

# # Install certificates for secure connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# # Copy the binary from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/worker .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go build`
CMD ["./main"]
CMD ["./worker"]