package models

type User struct {
    Email    string   `json:"email"`
    Limit    int      `json:"limit"`
    ActiveIPs []string `json:"active_ips"`
}
