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

type UserGameResponse struct {
	UserID        int                    `json:"userId"`
	GameErogsID   int                    `json:"gameErogsId"`
	Status        int                    `json:"status"`
	WishListMark  bool                   `json:"wishListMark"`
	BlackListMark bool                   `json:"blackListMark"`
	StartDate     *time.Time             `json:"startDate,omitempty"`
	FinishedDate  *time.Time             `json:"finishedDate,omitempty"`
	CreatedAt     time.Time              `json:"createdAt"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	GameErogs     *UserGameErogsResponse `json:"gameErogs,omitempty"`
}

type UserGameErogsResponse struct {
	ID           int                         `json:"id"`
	BrandErogsID int                         `json:"brandErogsId"`
	Name         string                      `json:"name"`
	Image        string                      `json:"image"`
	CreatedAt    time.Time                   `json:"createdAt"`
	UpdatedAt    time.Time                   `json:"updatedAt"`
	BrandErogs   *UserGameBrandErogsResponse `json:"brandErogs,omitempty"`
}

type UserGameBrandErogsResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Disband   bool      `json:"disband"`
	GameCount int       `json:"gameCount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RegisterResponse struct {
	DiscordID string `json:"discord_id"`
	UserName  string `json:"user_name"`
}

type RegisterLookupResponse struct {
	DiscordID string `json:"discord_id"`
}
