package models

type BlockedIP struct {
    IP       string `json:"ip"`
    BanTime  int    `json:"ban_time"`
    BannedAt int64  `json:"banned_at"`
}
