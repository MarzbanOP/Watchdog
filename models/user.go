package models

import (
	"time"
)

type User struct {
	Email     string    `json:"email"`
    Limit    int      `json:"limit"`
	ActiveIPs []string  `json:"active_ips"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}