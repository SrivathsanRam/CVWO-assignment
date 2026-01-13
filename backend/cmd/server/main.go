package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/SrivathsanRam/CVWO_project/backend/internal/database"
    "github.com/SrivathsanRam/CVWO_project/backend/internal/router"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file (only in development)
    _ = godotenv.Load(".env")

    // Initialize database connection
    if _, err := database.GetDB(); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.CloseDB()

    r := router.Setup()

    // Use PORT from environment (Render sets this automatically)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    fmt.Printf("Server starting on port %s\n", port)
    log.Fatalln(http.ListenAndServe(":"+port, r))
}