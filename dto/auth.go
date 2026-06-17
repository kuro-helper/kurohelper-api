package dto

import "time"

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// SessionUser is the logged-in user snapshot stored in session.
type SessionUser struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DiscordID   string    `json:"discordId"`
	Avatar      string    `json:"avatar"`
	Description string    `json:"description"`
	Role        int       `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type LoginResponse struct {
	User SessionUser `json:"user"`
}
