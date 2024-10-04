package main

import (
	"BigDataForge/internal/routes"
	"BigDataForge/internal/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set up Redis connection
	redisClient := storage.NewRedisClient()

	// Set up Gin router
	router := gin.Default()

	// Initialize routes
	routes.SetupRoutes(router, redisClient)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
