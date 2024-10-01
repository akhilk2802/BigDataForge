package main

import (
	"BigDataForge/config"
	"BigDataForge/routes"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize routes
	router := routes.SetupRouter()

	// Start the server
	router.Run(":8080")
}
