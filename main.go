package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kart-challenge/routes"
)

const defaultPort = "8080"

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	} else {
		log.Println("Loaded .env file")
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
		log.Printf("PORT not set, using default: %s", defaultPort)
	}
	return port
}

func main() {
	fmt.Println("kart challenge!")
	router := gin.Default()
	routes.InitializeRoutes(router)

	port := getPort()
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
