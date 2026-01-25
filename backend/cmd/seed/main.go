package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/SrivathsanRam/CVWO_project/backend/internal/database"
    "github.com/joho/godotenv"
)

func main() {
    // Load .env file
    envPaths := []string{".env", "../.env", "../../.env"}
    for _, path := range envPaths {
        if err := godotenv.Load(path); err == nil {
            log.Printf("Loaded .env from: %s", path)
            break
        }
    }

    // Connect to database
    db, err := database.GetDB()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.CloseDB()

    // Find seed.sql file
    seedPaths := []string{
        "migrations/seed.sql",
        "../migrations/seed.sql",
        "../../migrations/seed.sql",
    }

    var seedSQL []byte
    for _, path := range seedPaths {
        seedSQL, err = os.ReadFile(path)
        if err == nil {
            absPath, _ := filepath.Abs(path)
            log.Printf("Found seed.sql at: %s", absPath)
            break
        }
    }

    if seedSQL == nil {
        log.Fatal("Could not find seed.sql file")
    }

    // Execute seed SQL
    log.Println("Seeding database...")
    _, err = db.Exec(string(seedSQL))
    if err != nil {
        log.Fatalf("Failed to seed database: %v", err)
    }

    log.Println("Database seeded successfully!")
}