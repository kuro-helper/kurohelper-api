package dto

import "time"

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

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

type MeUserResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DiscordID   string    `json:"discordId"`
	Avatar      string    `json:"avatar"`
	Description string    `json:"description"`
	Role        int       `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type MeResponse struct {
	User *MeUserResponse `json:"user"`
}
