package dto

import "time"

type RegisterRequest struct {
	RegisterID string `json:"register_id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserHasPlayedResponse struct {
	UserID         string     `json:"userId"`
	GameErogsID    int        `json:"gameErogsId"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	GameID         int        `json:"gameId"`
	BrandID        int        `json:"brandId"`
	GameName       string     `json:"gameName"`
	GameImage      string     `json:"gameimage"`
	BrandName      string     `json:"brandName"`
	Disband        bool       `json:"disband"`
	BrandGameCount int        `json:"brandGameCount"`
}

type UserInWishResponse struct {
	UserID         string    `json:"userId"`
	GameErogsID    int       `json:"gameErogsId"`
	CreatedAt      time.Time `json:"createdAt"`
	GameID         int       `json:"gameId"`
	BrandID        int       `json:"brandId"`
	GameName       string    `json:"gameName"`
	GameImage      string    `json:"gameimage"`
	BrandName      string    `json:"brandName"`
	Disband        bool      `json:"disband"`
	BrandGameCount int       `json:"brandGameCount"`
}

type RegisterResponse struct {
	DiscordID string `json:"discord_id"`
	UserName  string `json:"user_name"`
}

type RegisterLookupResponse struct {
	DiscordID string `json:"discord_id"`
}
