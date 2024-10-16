package handlers

import (
	"encoding/json"
	"os"
	"time"
	"watchdog/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// APIAddUserRedis - Handler to add user in Redis
func APIAddUserRedis(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	// Set user in Redis
	if err := rdb.Set(ctx, user.Email, user.Limit, 0).Err(); err != nil {
		return c.Status(500).SendString("Failed to add user to Redis")
	}


	return c.Status(201).JSON(user)
}

// APIDeleteUserRedis - Handler to delete user in Redis

func APIDeleteUserRedis(c *fiber.Ctx) error {
	email := c.Params("email")

	// Delete user from Redis
	if err := rdb.Del(ctx, email).Err(); err != nil {
		return c.Status(500).SendString("Failed to delete user from Redis")
	}
	return c.Status(204).SendString("")
}

// APIAddUserJSON - Handler to add user in JSON
func APIAddUserJSON(c *fiber.Ctx) error {
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


	return c.Status(201).JSON(newUser)
}

// APIDeleteUserJSON - Handler to delete user in JSON
func APIDeleteUserJSON(c *fiber.Ctx) error {
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

// APIAddUserSQLite - Handler to add user in SQLite
func APIAddUserSQLite(c *fiber.Ctx, db *gorm.DB) error {
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(400).SendString("Invalid input")
	}

	if err := db.Create(&newUser).Error; err != nil {
		return c.Status(500).SendString("Failed to add user to SQLite")
	}

	return c.Status(201).JSON(newUser)
}

// APIDeleteUserSQLite - Handler to delete user in SQLite
func APIDeleteUserSQLite(c *fiber.Ctx, db *gorm.DB) error {
	email := c.Params("email")

	if err := db.Where("email = ?", email).Delete(&models.User{}).Error; err != nil {
		return c.Status(500).SendString("Failed to delete user from SQLite")
	}

	return c.Status(204).SendString("")
}


// APIBlockIPRedis - Handler to block an IP in Redis
func APIBlockIPRedis(c *fiber.Ctx) error {
	ip := c.Params("ip")
	banTime := 5 // Ban time in minutes

	// Set blocked IP in Redis
	if err := rdb.Set(ctx, ip, banTime, 0).Err(); err != nil {
		return c.Status(500).SendString("Failed to block IP")
	}

	return c.Status(200).SendString("IP blocked successfully")
}

// APIBlockIPJSON - Handler to block an IP in JSON
func APIBlockIPJSON(c *fiber.Ctx) error {
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

// APIBlockIPSQLite - Handler to block an IP in SQLite
func APIBlockIPSQLite(c *fiber.Ctx, db *gorm.DB) error {
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

// APIUnblockIPRedis - Handler to unblock an IP in Redis
func APIUnblockIPRedis(c *fiber.Ctx) error {
	ip := c.Params("ip")

	// Delete IP from Redis
	if err := rdb.Del(ctx, ip).Err(); err != nil {
		return c.Status(500).SendString("Failed to unblock IP")
	}

	return c.Status(200).SendString("IP unblocked successfully")
}

// APIUnblockIPJSON - Handler to unblock an IP in JSON
func APIUnblockIPJSON(c *fiber.Ctx) error {
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

// APIUnblockIPSQLite - Handler to unblock an IP in SQLite
func APIUnblockIPSQLite(c *fiber.Ctx, db *gorm.DB) error {
	ip := c.Params("ip")

	if err := db.Where("ip = ?", ip).Delete(&models.BlockedIP{}).Error; err != nil {
		return c.Status(500).SendString("Failed to unblock IP")
	}

	return c.Status(200).SendString("IP unblocked successfully")
}
