package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"watchdog/models"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	ctx = context.Background()
	rdb *redis.Client
	mu  sync.Mutex
	db  *gorm.DB // Global database connection
)

// Initialize Redis client
func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379", // Assuming Redis is running in a Docker container
	})
}

// InitSQLite - Initialize SQLite DB
func InitSQLite() (*gorm.DB, error) {
	// Placeholder for actual SQLite initialization
	// Replace with actual SQLite initialization code
	return nil, nil
}

// GetUserJSON function to retrieve a user by email from the JSON file
func GetUserJSON(email string, user *models.User) error {
	mu.Lock()
	defer mu.Unlock()

	// Read current users from users.json
	data, err := os.ReadFile("storage/users.json")
	if err != nil {
		return fmt.Errorf("failed to read users: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to parse users JSON: %w", err)
	}

	// Search for the user by email
	for _, u := range users {
		if u.Email == email {
			*user = u  // Copy the found user to the provided user pointer
			return nil // Successfully found and returned the user
		}
	}

	return fmt.Errorf("user with email %s not found", email) // User not found
}

// GetAllUserJSON retrieves all users from the JSON file
func GetAllUserJSON() ([]models.User, error) {
	mu.Lock()
	defer mu.Unlock()

	// Read current users from users.json
	data, err := os.ReadFile("storage/users.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read users: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users JSON: %w", err)
	}
	return users, nil
}

// AddUserRedis adds a user to Redis and manages their IPs
func AddUserRedis(user *models.User, newIP string) error {
	// Fetch existing user data from Redis
	existingData, err := rdb.Get(ctx, user.Email).Result()
	if err == redis.Nil {
		// If the user does not exist, create a new one
		user.ActiveIPs = []string{} // Initialize if creating new user
	} else if err != nil {
		return fmt.Errorf("failed to get user from Redis: %v", err)
	} else {
		// Unmarshal existing user data
		if err := json.Unmarshal([]byte(existingData), user); err != nil {
			return fmt.Errorf("failed to deserialize existing user: %v", err)
		}
	}

	// Check if the new IP is already in the ActiveIPs list
	for _, ip := range user.ActiveIPs {
		if ip == newIP {
			fmt.Printf("IP %s already exists for user %s\n", newIP, user.Email)
			return nil // Exit if the IP already exists
		}
	}

	// Add the new IP to ActiveIPs
	user.ActiveIPs = append(user.ActiveIPs, newIP)

	// Serialize updated user struct to JSON
	userData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to serialize user: %v", err)
	}

	// Store serialized user data in Redis
	if err := rdb.Set(ctx, user.Email, userData, 0).Err(); err != nil {
		return fmt.Errorf("failed to add/update user in Redis: %v", err)
	}

	return nil
}

// AddUserJSON adds a user to the JSON file and manages their IPs
func AddUserJSON(newUser *models.User, newIP string) error {
	mu.Lock()
	defer mu.Unlock()

	// Read current users from users.json
	data, err := os.ReadFile("storage/users.json")
	if err != nil {
		return fmt.Errorf("failed to read users: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to parse users JSON: %w", err)
	}

	// Check for existing user
	for i, user := range users {
		if user.Email == newUser.Email {
			fmt.Printf("User %s found, updating IPs...\n", user.Email)
			// Check if the new IP is already in the ActiveIPs list
			for _, ip := range user.ActiveIPs {
				if ip == newIP {
					fmt.Printf("IP %s already exists for user %s\n", newIP, newUser.Email)
					return nil // Exit if the IP already exists
				}
			}
			// Add the new IP to ActiveIPs
			users[i].ActiveIPs = append(users[i].ActiveIPs, newIP)
			// Update the updated_at timestamp
			users[i].UpdatedAt = time.Now()
			break
		}
	}

	// If the user doesn't exist, create a new user
	newUser.ActiveIPs = []string{newIP}
	newUser.CreatedAt = time.Now() // Set the created_at timestamp
	newUser.UpdatedAt = time.Now() // Set the updated_at timestamp
	users = append(users, *newUser)

	// Write updated users back to JSON file
	updatedData, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal updated users: %w", err)
	}

	if err := os.WriteFile("storage/users.json", updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write users to JSON: %w", err)
	}

	return nil
}

// AddUserSQLite adds a user to SQLite and manages their IPs
func AddUserSQLite(newUser *models.User, newIP string) error {
	mu.Lock()
	defer mu.Unlock()

	// Read current users from the database
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return fmt.Errorf("failed to read users from SQLite: %w", err)
	}

	// Check if the user already exists and update IPs
	for i, user := range users {
		if user.Email == newUser.Email {
			// Check if the new IP is already in the ActiveIPs list
			for _, ip := range user.ActiveIPs {
				if ip == newIP {
					fmt.Printf("IP %s already exists for user %s\n", newIP, newUser.Email)
					return nil // Exit if the IP already exists
				}
			}
			// Add the new IP to ActiveIPs
			users[i].ActiveIPs = append(users[i].ActiveIPs, newIP)
			// Update timestamps
			users[i].UpdatedAt = time.Now()
			// Update the user with new details
			if err := db.Save(&users[i]).Error; err != nil {
				return fmt.Errorf("failed to update user in SQLite: %w", err)
			}
			break
		}
	}

	// If the user doesn't exist, create a new user with the new IP
	newUser.ActiveIPs = []string{newIP}
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	if err := db.Create(newUser).Error; err != nil {
		return fmt.Errorf("failed to add user to SQLite: %w", err)
	}

	return nil
}
