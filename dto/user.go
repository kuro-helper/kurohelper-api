package dto

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserHasPlayed struct {
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

type UserInWish struct {
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

type RegisterRequest struct {
	RegisterID string `json:"register_id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
}

type RegisterResponse struct {
	DiscordID string `json:"discord_id"`
	UserName  string `json:"user_name"`
}

// RegisterLookupData 註冊邀請快取查詢結果（僅回傳 Discord ID）
type RegisterLookupData struct {
	DiscordID string `json:"discord_id"`
}
