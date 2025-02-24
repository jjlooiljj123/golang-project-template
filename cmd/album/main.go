package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"boilerplate/app/infrastructure/config"
	"boilerplate/app/infrastructure/httpclient"
	"boilerplate/app/infrastructure/httpclient/jsonpost"
	"boilerplate/app/infrastructure/redis"
	mysqlRepo "boilerplate/app/infrastructure/repositories/mysql"
	restcontroller "boilerplate/app/presentation/rest/album"
	"boilerplate/app/presentation/rest/router"
	albumservice "boilerplate/app/usecase/album"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Open MySQL connection
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", config.AppCfg.MySQLUser, config.AppCfg.MySQLPassword, config.AppCfg.MySQLHost, config.AppCfg.MySQLDatabase)
	db, err := mysqlRepo.OpenMySQLConnection(connectionString)
	if err != nil {
		log.Fatalf("Failed to open MySQL connection: %v", err)
	}

	// Initialize Redis cache
	redisCache, err := redis.NewRedisCache(config.AppCfg.RedisHost+":"+config.AppCfg.RedisPort, "", 0) // Adjust these according to your setup
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}

	// Initialize Repository layer
	albumRepo, err := mysqlRepo.NewAlbumRepository(db)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Initialize HTTP client
	httpClient := httpclient.NewClient()

	// Initialize third-party api service
	jsonPostHTTPClient := jsonpost.NewHttpJsonPost(httpClient, &config.AppCfg)

	// Initialize Usecase layer
	albumService := albumservice.NewService(albumRepo, redisCache, config.AppCfg.CacheDuration, jsonPostHTTPClient)

	// Initialize Controller layer
	restController := restcontroller.NewController(albumService)

	// set up routers
	r := gin.Default()
	router.SetupRoutes(r, restController, &config.AppCfg)

	// Start the server
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
