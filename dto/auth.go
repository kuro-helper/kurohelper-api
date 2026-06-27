package dto

import "time"

type LoginRequest struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User AuthUser `json:"user"`
}

type MeResponse struct {
	User *AuthUser `json:"user"`
}

type AuthUser struct {
	ID          int       `json:"id"`
	UserName    string    `json:"userName"`
	NickName    string    `json:"nickName"`
	DiscordID   string    `json:"discordId"`
	Avatar      string    `json:"avatar"`
	Description string    `json:"description"`
	Role        int       `json:"role"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
