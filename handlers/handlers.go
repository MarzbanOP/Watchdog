package handlers

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
	"watchdog/models"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


var (
	ctx                = context.Background()
	rdb                *redis.Client
	mu                 sync.Mutex // Mutex for JSON file operations
	userDeleteDelay, _ = strconv.Atoi(os.Getenv("USER_DELETE_DELAY")) // Read from .env
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

// AddUserRedis - Handler to add user in Redis
func AddUserRedis(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	// Set user in Redis
	if err := rdb.Set(ctx, user.Email, user.Limit, 0).Err(); err != nil {
		return c.Status(500).SendString("Failed to add user to Redis")
	}

	// Schedule user deletion after USER_DELETE_DELAY seconds
	go deleteUserAfterDuration(c,user.Email, "redis")

	return c.Status(201).JSON(user)
}

// DeleteUserRedis - Handler to delete user in Redis

func DeleteUserRedis(c *fiber.Ctx) error {
	email := c.Params("email")

	// Delete user from Redis
	if err := rdb.Del(ctx, email).Err(); err != nil {
		return c.Status(500).SendString("Failed to delete user from Redis")
	}
	return c.Status(204).SendString("")
}

// AddUserJSON - Handler to add user in JSON
func AddUserJSON(c *fiber.Ctx) error {
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	mu.Lock()
	defer mu.Unlock()

	// Read current users from users.json
	data, err := os.ReadFile("storage/users.json")
	if err != nil {
		return c.Status(500).SendString("Failed to read users")
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return c.Status(500).SendString("Failed to parse users JSON")
	}

	// Append new user
	users = append(users, newUser)

	// Write updated data back to file
	updatedData, err := json.Marshal(users)
	if err != nil {
		return c.Status(500).SendString("Failed to marshal updated users")
	}

	if err := os.WriteFile("storage/users.json", updatedData, 0644); err != nil {
		return c.Status(500).SendString("Failed to write updated users")
	}

	// Schedule user deletion after USER_DELETE_DELAY seconds
	go deleteUserAfterDuration(c,newUser.Email, "json")

	return c.Status(201).JSON(newUser)
}

// DeleteUserJSON - Handler to delete user in JSON
func DeleteUserJSON(c *fiber.Ctx) error {
	email := c.Params("email")

	mu.Lock()
	defer mu.Unlock()

	// Read current users from users.json
	data, err := os.ReadFile("storage/users.json")
	if err != nil {
		return c.Status(500).SendString("Failed to read users")
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return c.Status(500).SendString("Failed to parse users JSON")
	}

	// Find and delete the user
	for i, user := range users {
		if user.Email == email {
			users = append(users[:i], users[i+1:]...) // Remove user
			break
		}
	}

	// Write updated data back to file
	updatedData, err := json.Marshal(users)
	if err != nil {
		return c.Status(500).SendString("Failed to marshal updated users")
	}

	if err := os.WriteFile("storage/users.json", updatedData, 0644); err != nil {
		return c.Status(500).SendString("Failed to write updated users")
	}

	return c.Status(204).SendString("")
}

// AddUserSQLite - Handler to add user in SQLite
func AddUserSQLite(c *fiber.Ctx, db *gorm.DB) error {
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	if err := db.Create(&newUser).Error; err != nil {
		return c.Status(500).SendString("Failed to add user to SQLite")
	}

	// Schedule user deletion after USER_DELETE_DELAY seconds
	go deleteUserAfterDuration(c,newUser.Email, "sqlite")

	return c.Status(201).JSON(newUser)
}

// DeleteUserSQLite - Handler to delete user in SQLite
func DeleteUserSQLite(c *fiber.Ctx, db *gorm.DB) error {
	email := c.Params("email")

	if err := db.Where("email = ?", email).Delete(&models.User{}).Error; err != nil {
		return c.Status(500).SendString("Failed to delete user from SQLite")
	}

	return c.Status(204).SendString("")
}


// deleteUserAfterDuration deletes the user after a specified duration
func deleteUserAfterDuration(c *fiber.Ctx, email, storageType string) {
	time.Sleep(time.Duration(userDeleteDelay) * time.Second)

	// Call the appropriate delete function based on storage type
	switch storageType {
	case "redis":
		_ = DeleteUserRedis(c) // Error ignored for simplicity

	case "json":
		_ = DeleteUserJSON(c) // Error ignored for simplicity

	case "sqlite":
		_ = DeleteUserSQLite(c, nil) // Ensure to pass a valid db instance
	}
}


// BlockIPRedis - Handler to block an IP in Redis
func BlockIPRedis(c *fiber.Ctx) error {
	ip := c.Params("ip")
	banTime := 5 // Ban time in minutes

	// Set blocked IP in Redis
	if err := rdb.Set(ctx, ip, banTime, 0).Err(); err != nil {
		return c.Status(500).SendString("Failed to block IP")
	}

	return c.Status(200).SendString("IP blocked successfully")
}

// BlockIPJSON - Handler to block an IP in JSON
func BlockIPJSON(c *fiber.Ctx) error {
	ip := c.Params("ip")
	banTime := 10 // Ban time in minutes

	mu.Lock()
	defer mu.Unlock()

	// Read existing blocked IPs from blocked_ips.json
	data, err := os.ReadFile("storage/blocked_ips.json")
	if err != nil {
		return c.Status(500).SendString("Failed to read blocked IPs")
	}

	var blockedIPs []models.BlockedIP
	if err := json.Unmarshal(data, &blockedIPs); err != nil {
		return c.Status(500).SendString("Failed to parse blocked IPs JSON")
	}

	// Add new IP to blocked list
	newBlockedIP := models.BlockedIP{
		IP:       ip,
		BanTime:  banTime,
		BannedAt: time.Now().Unix(),
	}
	blockedIPs = append(blockedIPs, newBlockedIP)

	// Write updated data back to file
	updatedData, err := json.Marshal(blockedIPs)
	if err != nil {
		return c.Status(500).SendString("Failed to marshal updated blocked IPs")
	}

	if err := os.WriteFile("storage/blocked_ips.json", updatedData, 0644); err != nil {
		return c.Status(500).SendString("Failed to write updated blocked IPs")
	}

	return c.Status(200).SendString("IP blocked successfully")
}

// BlockIPSQLite - Handler to block an IP in SQLite
func BlockIPSQLite(c *fiber.Ctx, db *gorm.DB) error {
	ip := c.Params("ip")
	banTime := 10 // Ban time in minutes

	blockedIP := models.BlockedIP{
		IP:       ip,
		BanTime:  banTime,
		BannedAt: time.Now().Unix(),
	}

	if err := db.Create(&blockedIP).Error; err != nil {
		return c.Status(500).SendString("Failed to block IP")
	}

	return c.Status(200).SendString("IP blocked successfully")
}

// UnblockIPRedis - Handler to unblock an IP in Redis
func UnblockIPRedis(c *fiber.Ctx) error {
	ip := c.Params("ip")

	// Delete IP from Redis
	if err := rdb.Del(ctx, ip).Err(); err != nil {
		return c.Status(500).SendString("Failed to unblock IP")
	}

	return c.Status(200).SendString("IP unblocked successfully")
}

// UnblockIPJSON - Handler to unblock an IP in JSON
func UnblockIPJSON(c *fiber.Ctx) error {
	ip := c.Params("ip")

	mu.Lock()
	defer mu.Unlock()

	// Read existing blocked IPs from blocked_ips.json
	data, err := os.ReadFile("storage/blocked_ips.json")
	if err != nil {
		return c.Status(500).SendString("Failed to read blocked IPs")
	}

	var blockedIPs []models.BlockedIP
	if err := json.Unmarshal(data, &blockedIPs); err != nil {
		return c.Status(500).SendString("Failed to parse blocked IPs JSON")
	}

	// Find and delete the IP
	for i, blockedIP := range blockedIPs {
		if blockedIP.IP == ip {
			blockedIPs = append(blockedIPs[:i], blockedIPs[i+1:]...) // Remove IP
			break
		}
	}

	// Write updated data back to file
	updatedData, err := json.Marshal(blockedIPs)
	if err != nil {
		return c.Status(500).SendString("Failed to marshal updated blocked IPs")
	}

	if err := os.WriteFile("storage/blocked_ips.json", updatedData, 0644); err != nil {
		return c.Status(500).SendString("Failed to write updated blocked IPs")
	}

	return c.Status(200).SendString("IP unblocked successfully")
}

// UnblockIPSQLite - Handler to unblock an IP in SQLite
func UnblockIPSQLite(c *fiber.Ctx, db *gorm.DB) error {
	ip := c.Params("ip")

	if err := db.Where("ip = ?", ip).Delete(&models.BlockedIP{}).Error; err != nil {
		return c.Status(500).SendString("Failed to unblock IP")
	}

	return c.Status(200).SendString("IP unblocked successfully")
}


