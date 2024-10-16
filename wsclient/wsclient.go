package wsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"watchdog/handlers"
	"watchdog/models"

	"github.com/gorilla/websocket"
)

// Structure for the token response
type TokenResponse struct {
    AccessToken string `json:"access_token"`
}

// Function to get the token by sending admin credentials
func GetToken() (string, error) {
    address := os.Getenv("ADDRESS")
    port := os.Getenv("PORT_ADDRESS")
    // Construct the token URL
    tokenURL := fmt.Sprintf("http://%s:%s/api/admin/token", address, port)

    data := url.Values{}
    data.Set("grant_type", "password")
    data.Set("username", os.Getenv("P_USER")) // Load from .env
    data.Set("password", os.Getenv("P_PASS")) // Load from .env
    data.Set("scope", "")
    data.Set("client_id", "")
    data.Set("client_secret", "")

    resp, err := http.PostForm(tokenURL, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to authenticate: %s", resp.Status)
    }

    var tokenResponse TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return "", err
    }

    return tokenResponse.AccessToken, nil
}

// Function to connect to WebSocket with token
func ConnectToWebSocket(token string) {
    // Check if SSL is enabled by checking if SSL is set to "true"
    ssl := os.Getenv("SSL") == "true"
    address := os.Getenv("ADDRESS")
    port := os.Getenv("PORT_ADDRESS")
    serverURL := fmt.Sprintf("%s:%s", address, port) // Store URL in .env
    interval := os.Getenv("LOG_INTERVAL") // Log interval from .env
    if interval == "" {
        interval = "5"
    }

    // Construct the WebSocket URL based on SSL flag
    wsURL := fmt.Sprintf("%s://%s/api/core/logs?interval=%s",
        map[bool]string{true: "wss", false: "ws"}[ssl], serverURL, interval)

    headers := http.Header{}
    headers.Add("Authorization", "Bearer "+token)

    log.Printf("Connecting to WebSocket at %s", wsURL)
    c, _, err := websocket.DefaultDialer.Dial(wsURL, headers)
    if err != nil {
        log.Printf("Connection error: %v", err)
        return
    }
    defer c.Close()

    // Send initial message
    err = c.WriteMessage(websocket.TextMessage, []byte("Hello, Server!"))
    if err != nil {
        log.Fatalf("Error sending message: %v", err)
        return
    }

    // Continuously read messages
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            return
        }
        parseMessage(string(message))
    }
}


// parseMessage extracts IP and email using regex
func parseMessage(message string) (string, string) {
    ipRegex := regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
    emailRegex := regexp.MustCompile(`email:\s*([^\s]+)`)

    ipMatch := ipRegex.FindStringSubmatch(message)
    emailMatch := emailRegex.FindStringSubmatch(message)

    if len(ipMatch) < 2 || len(emailMatch) < 2 {
        return "", ""
    }
    sendToStorage(ipMatch[1], emailMatch[1])
    return ipMatch[1], emailMatch[1]
}

// sendToStorage sends the extracted data, LIMIT, email, and ActiveIPs to the /api/user/add endpoint at port 4000
func sendToStorage(ip, email string) {
    storageType := os.Getenv("STORAGE_TYPE")
    if storageType == "" {
        log.Fatal("STORAGE_TYPE environment variable is not set")
    }

    limitStr := os.Getenv("MAX_ALLOW_USERS")
    limit := 0
    if limitStr != "" {
        if _, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil {
            log.Printf("Error parsing LIMIT: %v", err)
            return
        }
    }

    var user models.User

    // Retrieve existing user data from storage
    switch storageType {
    // case "redis":
    //     if err := handlers.GetUserRedis(email, &user); err != nil {
    //         log.Printf("Error retrieving user from Redis: %v", err)
    //     }

    case "json":
        if err := handlers.GetUserJSON(email, &user); err != nil {
            log.Printf("Error retrieving user from JSON: %v", err)
        }

    // case "sqlite":
    //     if err := handlers.GetUserSQLite(email, &user); err != nil {
    //         log.Printf("Error retrieving user from SQLite: %v", err)
    //     }

    default:
        log.Printf("Unknown storage type: %s", storageType)
        return
    }

    // Add new IP to the ActiveIPs slice if it doesn't already exist
    if !contains(user.ActiveIPs, ip) {
        user.ActiveIPs = append([]string{ip}, user.ActiveIPs...) // Add new IP at the beginning
    } else {
        log.Printf("IP %s is already in the user's active IPs.", ip)
    }

    user.Email = email
    user.Limit = limit // Update limit if necessary

    // Store updated user data back to storage
    switch storageType {
    // case "redis":
    //     if err := handlers.AddUserRedis(&user,ip); err != nil {
    //         log.Printf("Error storing user in Redis: %v", err)
    //     } else {
    //         log.Println("User successfully stored in Redis.")
    //     }

    case "json":
        if err := handlers.AddUserJSON(&user,ip); err != nil {
            log.Printf("Error storing user in JSON: %v", err)
        } else {
            log.Println("User successfully stored in JSON.")
        }

    // case "sqlite":
    //     if err := handlers.AddUserSQLite(db,&user); err != nil {
    //         log.Printf("Error storing user in SQLite: %v", err)
    //     } else {
    //         log.Println("User successfully stored in SQLite.")
    //     }

    default:
        log.Printf("Unknown storage type: %s", storageType)
    }

    // Optionally, you can marshal the user data to JSON after storage
    jsonData, err := json.Marshal(user)
    if err != nil {
        log.Printf("Error marshalling JSON: %v", err)
        return
    }

    // For demonstration, just logging the marshaled JSON data
    log.Printf("User data in JSON format: %s\n", jsonData)
}

// Helper function to check if an IP is already in the ActiveIPs slice
func contains(slice []string, ip string) bool {
    for _, item := range slice {
        if item == ip {
            return true
        }
    }
    return false
}