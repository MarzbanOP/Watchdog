package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

// Structure for the token response
type TokenResponse struct {
    AccessToken string `json:"access_token"`
}

// Function to get the token by sending admin credentials
func getToken() (string, error) {
    // API URL for token authentication
    tokenURL := "http://37.152.189.61:8000/api/admin/token"

    // Prepare the data for the form
    data := url.Values{}
    data.Set("grant_type", "password")
    data.Set("username", "admin")
    data.Set("password", "admin")
    data.Set("scope", "")
    data.Set("client_id", "")       // Set your client ID here if required
    data.Set("client_secret", "")   // Set your client secret here if required

    // Send a POST request to /api/admin/token to get the token
    resp, err := http.PostForm(tokenURL, data)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // Check if the request was successful
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("failed to authenticate: %s", resp.Status)
    }

    // Parse the response to extract the token
    var tokenResponse TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return "", err
    }

    return tokenResponse.AccessToken, nil
}

// Function to construct the WebSocket URL and connect with the Bearer token
func connectToWebSocket(token string) {
    // WebSocket server details
    ssl := os.Getenv("SSL") == "false"  // Check if SSL is enabled
    url := "37.152.189.61:8000"          // Replace with the appropriate server URL
    interval := "1"                     // Set your desired interval

    // Construct the WebSocket URL
    serverURL := fmt.Sprintf("%s://%s/api/core/logs?interval=%s",
        map[bool]string{true: "wss", false: "ws"}[ssl], url, interval)

    // Create a custom header with the Bearer token
    headers := http.Header{}
    headers.Add("Authorization", "Bearer "+token)

    // Log the server URL and headers for debugging
    log.Printf("Connecting to WebSocket at %s with headers: %v", serverURL, headers)

    // Dial WebSocket with custom headers
    c, _, err := websocket.DefaultDialer.Dial(serverURL, headers)
    if err != nil {
        log.Printf("Connection error: %v", err)
        return
    }
    defer c.Close()

    // Send a test message to the WebSocket server
    err = c.WriteMessage(websocket.TextMessage, []byte("Hello, Server!"))
    if err != nil {
        log.Fatalf("Error sending message: %v", err)
        return
    }

    // Wait for and log all responses
    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            return
        }
        logParsedMessage(string(message))
    }
}

// logParsedMessage sends extracted IP and email to the /api/user/add endpoint
func logParsedMessage(message string) {
    ip, email := parseMessage(message)
    if ip != "" && email != "" {
        sendToAPICallback(ip, email)
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

    return ipMatch[1], emailMatch[1]
}

// sendToAPICallback sends the extracted data, LIMIT, email, and ActiveIPs to the /api/user/add endpoint at port 4000
func sendToAPICallback(ip, email string) {
    limit := os.Getenv("MAX_ALLOW_USERS")

    apiURL := "http://localhost:4000/api/user/add" // API endpoint
    data := map[string]interface{}{
        "ip":        ip,
        "email":     email,
        "limit":     limit,
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Printf("Error marshalling JSON: %v", err)
        return
    }

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error creating HTTP request: %v", err)
        return
    }

    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error sending HTTP request: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        log.Printf("API callback failed: %v", resp.Status)
    } else {
        log.Println("Data successfully sent to /api/user/add.")
    }
}

func main() {
    // First, get the Bearer token by authenticating
    token, err := getToken()
    if err != nil {
        log.Fatalf("Error getting token: %v", err)
        return
    }

    // Now connect to the WebSocket server with the token
    for {
        connectToWebSocket(token)

        // If the connection is lost, reconnect after 5 seconds
        log.Println("Reconnecting in 5 seconds...")
        time.Sleep(5 * time.Second)
    }
}
