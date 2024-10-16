package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	"watchdog/handlers"
	"watchdog/wsclient"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// checkAndDeleteUsers checks users in the JSON file and deletes them if necessary
func checkAndDeleteUsers(queue chan<- string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userDeleteDelay , _ := strconv.Atoi(os.Getenv("USER_DELETE_DELAY"))
    users, err := handlers.GetAllUserJSON()
    if err != nil {
        fmt.Println("Error retrieving users:", err)
        return
    }

    currentTime := time.Now()
    for _, user := range users {
        // Calculate the time to delete based on UpdatedAt and userDeleteDelay 
        timeToDelete := user.UpdatedAt.Add(time.Duration(userDeleteDelay) * time.Second)
        if currentTime.After(timeToDelete) {
            // Send user email to the queue for deletion
            queue <- user.Email
            fmt.Printf("User %s is scheduled for deletion\n", user.Email)

            // Here, implement removal logic from the JSON file if necessary
        }
    }
}


func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	sleepDuration , _ := strconv.Atoi(os.Getenv("SLEEP_DURATION")) // New line to read sleep duration

	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		log.Fatal("STORAGE_TYPE environment variable is not set")
	}
	app := fiber.New()

	// WebSocket authentication and connection in a goroutine
	token, err := wsclient.GetToken()
	if err != nil {
		log.Fatalf("Error getting token: %v", err)
		return
	}

	go func() {
		for {
			wsclient.ConnectToWebSocket(token)
			log.Println("Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}()

	// Declare the queue channel with a buffer size of 5
	queue := make(chan string, 5)

	// Start a goroutine to handle user deletions
	go func() {
		for {
			checkAndDeleteUsers(queue) // Call the function that checks for user deletions
			time.Sleep(time.Duration(sleepDuration) * time.Second) // Sleep
		}
	}()
	// Initialize storage handlers based on environment variable
	switch storageType {
	case "redis":
		handlers.InitRedis()
		app.Post("/api/user/add", handlers.APIAddUserRedis)
		app.Delete("/api/user/delete/:email", handlers.APIDeleteUserRedis)
		app.Post("/api/ip/block/:ip", handlers.APIBlockIPRedis)
		app.Post("/api/ip/unblock/:ip", handlers.APIUnblockIPRedis)
	case "json":
		app.Post("/api/user/add", handlers.APIAddUserJSON)
		app.Delete("/api/user/delete/:email", handlers.APIDeleteUserJSON)
		app.Post("/api/ip/block/:ip", handlers.APIBlockIPJSON)
		app.Post("/api/ip/unblock/:ip", handlers.APIUnblockIPJSON)
	case "sqlite":
		db, err := handlers.InitSQLite()
		if err != nil {
			log.Fatal("Failed to connect to SQLite:", err)
		}
		app.Post("/api/user/add", func(c *fiber.Ctx) error {
			return handlers.APIAddUserSQLite(c, db)
		})
		app.Delete("/api/user/delete/:email", func(c *fiber.Ctx) error {
			return handlers.APIDeleteUserSQLite(c, db)
		})
		app.Post("/api/ip/block/:ip", func(c *fiber.Ctx) error {
			return handlers.APIBlockIPSQLite(c, db)
		})
		app.Post("/api/ip/unblock/:ip", func(c *fiber.Ctx) error {
			return handlers.APIUnblockIPSQLite(c, db)
		})
	default:
		log.Fatal("Invalid STORAGE_TYPE specified. Must be 'redis', 'json', or 'sqlite'.")
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "4000"
	}
	log.Fatal(app.Listen(":" + port))
}


