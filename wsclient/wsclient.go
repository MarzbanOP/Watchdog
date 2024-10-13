package wsclient

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"
    "regexp"
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

    // Set default values
    if address == "" {
        address = "127.0.0.1" // Default address
    }
    if port == "" {
        port = "8000" // Default port
    }

    // Construct the token URL
    tokenURL := fmt.Sprintf("http://%s:%s/api/admin/token", address, port)

    data := url.Values{}
    data.Set("grant_type", "password")
    data.Set("username", os.Getenv("ADMIN_USERNAME")) // Load from .env
    data.Set("password", os.Getenv("ADMIN_PASSWORD")) // Load from .env
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
    ssl := os.Getenv("SSL") == "false"
    serverURL := os.Getenv("WEBSOCKET_SERVER_URL") // Store URL in .env
    interval := os.Getenv("LOG_INTERVAL") // Log interval from .env
    if interval == "" {
        interval = "5"
    }
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

    err = c.WriteMessage(websocket.TextMessage, []byte("Hello, Server!"))
    if err != nil {
        log.Fatalf("Error sending message: %v", err)
        return
    }

    for {
        _, message, err := c.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            return
        }
        logParsedMessage(string(message))
    }
}

// logParsedMessage prints extracted IP and email
func logParsedMessage(message string) {
    ip, email := parseMessage(message)
    if ip != "" && email != "" {
        fmt.Printf("IP: %s, Email: %s\n", ip, email)
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
