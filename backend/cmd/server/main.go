package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SrivathsanRam/CVWO_project/backend/internal/database"
	"github.com/SrivathsanRam/CVWO_project/backend/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file from backend directory
	// Try current directory first, then parent directories
	envPaths := []string{".env", "../.env", "../../.env"}
	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env from: %s", path)
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database connection
	if _, err := database.GetDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseDB()

	r := router.Setup()
	fmt.Println("Listening on port 8000 at http://localhost:8000")

	log.Fatalln(http.ListenAndServe(":8000", r))
}
