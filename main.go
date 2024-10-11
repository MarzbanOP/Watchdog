package main

import (
    "github.com/gofiber/fiber/v2"
    "github.com/joho/godotenv"
    "log"
    "os"
    "ipwatchdog/handlers"
)

func main() {
    // Load environment variables
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Determine storage type
    storageType := os.Getenv("STORAGE_TYPE")
    if storageType == "" {
        log.Fatal("STORAGE_TYPE environment variable is not set")
    }

    // Create new Fiber instance
    app := fiber.New()

    // Initialize storage handlers based on environment variable
    switch storageType {
    case "redis":
        handlers.InitRedis()
        app.Post("/api/user/add", handlers.AddUserRedis)
        app.Delete("/api/user/delete/:email", handlers.DeleteUserRedis)
        app.Post("/api/ip/block/:ip", handlers.BlockIPRedis)
        app.Post("/api/ip/unblock/:ip", handlers.UnblockIPRedis)
    case "json":
        app.Post("/api/user/add", handlers.AddUserJSON)
        app.Delete("/api/user/delete/:email", handlers.DeleteUserJSON)
        app.Post("/api/ip/block/:ip", handlers.BlockIPJSON)
        app.Post("/api/ip/unblock/:ip", handlers.UnblockIPJSON)
    case "sqlite":
        // You would need to pass the db instance here, which could be initialized
        db, err := handlers.InitSQLite() // Assume this function initializes and returns the SQLite DB instance
        if err != nil {
            log.Fatal("Failed to connect to SQLite:", err)
        }
        app.Post("/api/user/add", func(c *fiber.Ctx) error {
            return handlers.AddUserSQLite(c, db)
        })
        app.Delete("/api/user/delete/:email", func(c *fiber.Ctx) error {
            return handlers.DeleteUserSQLite(c, db)
        })
        app.Post("/api/ip/block/:ip", func(c *fiber.Ctx) error {
            return handlers.BlockIPSQLite(c, db)
        })
        app.Post("/api/ip/unblock/:ip", func(c *fiber.Ctx) error {
            return handlers.UnblockIPSQLite(c, db)
        })
    default:
        log.Fatal("Invalid STORAGE_TYPE specified. Must be 'redis', 'json', or 'sqlite'.")
    }

    // Start server on defined port
    port := os.Getenv("API_PORT")
    if port == "" {
        port = "4000"
    }
    log.Fatal(app.Listen(":" + port))
}
