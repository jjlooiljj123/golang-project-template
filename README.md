# Golang Project Structure Template

This project showcases a Golang template demonstrating how to structure files for an HTTP application and an AWS SQS application. The HTTP application is related to albums, so some folder structures may include the 'album' naming convention.

What is included:

1. Consuming HTTP requests using the Gin router
2. Consuming AWS SQS messages
3. Dependency injection during application initialization
4. Sample middleware (extracting common headers, authentication, timeout, latency tracking)
5. Domain-driven design concepts (entities are used to pass data between different layers)
6. Using values obtained from environment variables
7. Sample code for making an HTTP call to a third-party API
8. Sample code for using Redis for caching
9. Sample code for interacting with MySQL
10. Sample code for interacting with AWS SQS
11. Dockerfile to containerize the applications
12. Docker Compose for containerizing dependencies (AWS SQS, Redis, MySQL) to run the app locally

## Project Structure

```bash
projectname/
├── cmd/                       # Main application entry points
│   ├── album/                 # Main HTTP application entry point
│   ├── worker/                # Main SQS application entry point
├── app/                       # All app logic folders entry point
│   ├── presentation/          # Entry point logic for HTTP (and other technologies like gRPC)
│   │   ├── rest/              # HTTP controllers entry point
│   │   │   ├── album/         # HTTP controllers for albums
│   │   │   ├── middleware/    # HTTP controllers middleware
│   │   │   ├── router/        # HTTP endpoints paths configuration
│   ├── domain/                # Entity object folder
│   │   ├── entity/            # Entity objects used to pass data between presentation, usecase, and infrastructure layers
│   │   ├── dto/               # DTOs for HTTP requests and responses
│   │   ├── errors/            # Custom error objects used in the repository
│   ├── usecase/               # Business logic folder
│   │   ├── album/             # Business logic for the HTTP application
│   │   ├── worker/            # Business logic for the SQS application
│   │   ├── interface/         # Interfaces for business logic, designed for dependency injection
│   ├── infrastructure/        # Entry point for folders interacting with external services or infrastructure
│   │   ├── repositories/      # Data storage folders (could contain MySQL, DynamoDB, etc.)
│   │   │   ├── interface/     # Interfaces for repository logic, designed for dependency injection
│   │   │   ├── mysql/         # MySQL logic
│   │   ├── redis/             # Redis logic
│   │   ├── httpclient/        # Entry point for interacting with external services using HTTP
│   │   │   ├── interface/     # Interfaces for third-party client logic, designed for dependency injection
│   │   │   ├── jsonpost/      # Sample third-party client interaction logic (making HTTP calls)
│   │   ├── sqs/               # Logic for consuming/sending SQS messages
│   │   ├── config/            # Config object for environment variables
├── scripts/                   # Contains all scripts (used for repo initialization, etc.)
├── resources/                 # Contains non-implementation-related items
├── .env                       # Environment variables
├── Dockerfile                 # Dockerfile
├── docker-compose.yaml        # Docker Compose file to build resources needed to run the app locally
├── .dockerignore              # Docker ignore file
├── .gitignore                 # Git ignore file
├── go.mod                     # Module dependencies
├── go.sum                     # Dependency checksums
└── README.md                  # This file
```

## Project Structure Diagrams

### HTTP Application

![http_architecture](./resources/diagrams/http_architecture.png?raw=true)

#### Entry Point [cmd/album/]

- Entry point of the HTTP application
- Loads environment variables
- Initializes services and dependencies (including dependency injection)
- Starts the HTTP server

#### Loading Environment Variables [app/infrastructure/config/]

- Defines the environment variables object
- Loads environment variables (from the <code>.env</code> file) into the defined object

#### Setting Up Routers [app/presentation/rest/router/]

- up HTTP routes
- Groups routes by prefix (<code>/v1</code>, <code>/v2</code>, ...)
- Adds middleware for respective routes

#### Middleware [app/presentation/rest/middleware/]

- Authentication, common header extractor, timeout, latency logger

#### Controller [app/presentation/rest/album/]

- Handles HTTP requests
- Captures HTTP request payloads (via DTO objects)
- Converts DTOs to entities before passing data to the usecase layer (business logic layer)
- Usecase layer in the controller is represented by interfaces.

<details>
<summary>Dependency Injection Example</summary>

```go
// --- In app/usecase/interface/portfolio_interface.go ---
type GetAlbumInterface interface {
	GetAllAlbums(ctx context.Context) ([]dto.Album, error)
	GetAlbumByID(ctx context.Context, id string) (dto.Album, error)
}

type CreateAlbumInterface interface {
	CreateAlbum(ctx context.Context, album entity.Album) (string, error)
}

type GetJSONPostInterface interface {
	GetFromThirdPartyAPI(ctx context.Context) ([]dto.Post, error)
}

type AlbumInterface interface {
	GetAlbumInterface
	CreateAlbumInterface
	GetJSONPostInterface
}

// --- In app/presentation/rest/album/controller.go ---
type Controller struct {
	albumService albumservice.AlbumInterface
}

// --- In cmd/album/main.go ---
// Initialize Usecase layer
albumService := albumservice.NewService(albumRepo, redisCache, config.AppCfg.CacheDuration, jsonPostHTTPClient)

// Initialize Controller layer
restController := restcontroller.NewController(albumService)
```

</details>

#### Usecase (Business Logic) Layer [app/usecase/album/]

- Handles all business logic
- Data is passed in and out via entity objects
- Infrastructure layers (repository, Redis, etc.) in the usecase layer are represented by interfaces

<details>
<summary>Dependency Injection Example</summary>

```go
// --- In app/infrastructure/repositories/interface/repo_interface.go ---
type RepositoryInterface interface {
	GetAlbums(ctx context.Context) ([]entity.Album, error)
	CreateAlbum(ctx context.Context, album entity.Album) (string, error)
	GetAlbumByID(ctx context.Context, id string) (entity.Album, error)
}


type AlbumInterface interface {
	GetAlbumInterface
	CreateAlbumInterface
	GetJSONPostInterface
}

// --- In app/usecase/album/service.go ---
type Service struct {
	albumRepo albumsRepositories.RepositoryInterface\
}

// --- In cmd/album/main.go ---
// Initialize Repository layer
albumRepo, err := mysqlRepo.NewAlbumRepository(db)

// Initialize Usecase layer
albumService := albumservice.NewService(albumRepo, redisCache, config.AppCfg.CacheDuration, jsonPostHTTPClient)
```

</details>

#### Repository Layer [app/infrastructure/repositories/]

- Handles all database-related logic
- Data is passed in and out via entity objects
- Entity objects are parsed into a format that can interact with the database

#### Redis Layer [app/infrastructure/redis/]

- Handles all Redis-related logic

#### Http Client Layer [app/infrastructure/httpclient/]

- Handles all HTTP calls (to external services) logic
- Data is passed in and out via entity objects
- Entity objects are parsed into requests for third-party services, and responses are parsed back into entity objects

### SQS Application

![sqs_architecture](./resources/diagrams/sqs_architecture.png?raw=true)

#### Entry Point [cmd/worker/]

- Entry point of the SQS application
- Loads environment variables
- Initializes services and dependencies (including dependency injection)
- Starts worker(s) consuming SQS messages

#### Loading Environment Variables [app/infrastructure/config/]

- Defines the environment variables object
- Loads environment variables (from the <code>.env</code> file) into the defined object

#### Setting Up Worker [app/usecase/worker/worker.go/]

- Initializes workers for SQS consumption

#### Message Handler [app/infrastructure/sqs/album_processor.go/]

- Logic injected into the SQS worker
- Consumes SQS messages, extracts the message body, and passes it to the usecase layer
- Deletes the message if no error is returned from the usecase layer processing

## Prerequisites

- Golang 1.23 or above
- Docker

## Installation

1. Clone the repository:

```bash
git clone https://github.com/jjlooiljj123/golang-project-template
```

2. Navigate to the project directory:

```bash
cd boilerplate
```

3. Install dependencies:

```bash
go mod tidy
go mod vendor
```

4. Build and run dependencies:

```bash
docker-compose up --build
```

4. Set up tables in local MySQL:

```bash
go run ./script/prepare_mysql_table/prepare_mysql_table.go
go run ./script/insert_dummy_data/insert_dummy_data.go
```

## Usage

Before running the app locally, update the following environment variables in the <code>.env</code> file:

- <code>MYSQL_HOST=localhost</code>
- <code>REDIS_HOST=localhost</code>
- <code>SQS_QUEUE_URL=http://localhost:4566/000000000000/album</code>
- <code>AWS_SQS_HOST=http://localhost:4566</code>

host values in the existing <code>.env</code> file are not set to <code>localhost</code> because they refer to the containers defined in the <code>docker-compose.yaml</code> file.

Run the HTTP application:

```bash
go run ./cmd/album/main.go
```

Run the SQS worker application:

```bash
go run ./cmd/worker/worker.go
```

Send a sample SQS message:

```bash
./script/send_sqs/send_sqs.sh ./script/send_sqs/sample_sqs.json
```

## Resources

#### API Collections

Please find all the cUrl here: <code>./resources/api_curl</code>
