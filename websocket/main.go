package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "regexp"
    "os"
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
    tokenURL := "http://5.75.200.248:8000/api/admin/token"

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
    url := "5.75.200.248:8000"          // Replace with the appropriate server URL
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

// logParsedMessage prints the extracted IP address and email from the message
func logParsedMessage(message string) {
    // Parse the message for IP and email
    ip, email := parseMessage(message)
    if ip != "" && email != "" { // Only log if both IP and email are found
        fmt.Printf("IP: %s, Email: %s\n", ip, email)
    }
}

// parseMessage extracts the IP address and email from the message using regex
func parseMessage(message string) (string, string) {
    // Regular expression to capture the IP address (before the colon) and the email (after 'email:')
    ipRegex := regexp.MustCompile(`([0-9]+\.[0-9]+\.[0-9]+\.[0-9]+)`)
    emailRegex := regexp.MustCompile(`email:\s*([^\s]+)`)

    // Extract IP address
    ipMatch := ipRegex.FindStringSubmatch(message)
    if len(ipMatch) < 2 {
        return "", "" // No IP found
    }
    ip := ipMatch[1]

    // Extract email
    emailMatch := emailRegex.FindStringSubmatch(message)
    if len(emailMatch) < 2 {
        return "", "" // No email found
    }
    email := emailMatch[1]

    return ip, email
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
