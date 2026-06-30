package dto

import "time"

type RegisterRequest struct {
	RegisterID string `json:"registerId"`
	UserName   string `json:"userName"`
	Password   string `json:"password"`
}

type UserResponse struct {
	ID              int       `json:"id"`
	NickName        string    `json:"nickName"`
	DiscordID       string    `json:"discordId"`
	Avatar          string    `json:"avatar"`
	PrivateGameData bool      `json:"privateGameData"`
	Role            int       `json:"role"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type UserProfileResponse struct {
	ID              int       `json:"id"`
	UserName        string    `json:"userName,omitempty"`
	NickName        string    `json:"nickName"`
	DiscordID       string    `json:"discordId"`
	Avatar          string    `json:"avatar"`
	Description     string    `json:"description"`
	PrivateGameData bool      `json:"privateGameData"`
	Role            int       `json:"role"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type GetUserGameResponse struct {
	User  UserProfileResponse  `json:"user"`
	Games []UserGameResponse   `json:"games"`
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
	DiscordID string `json:"discordId"`
	UserName  string `json:"userName"`
}

type RegisterLookupResponse struct {
	DiscordID string `json:"discordId"`
}

type UpdateUserRequest struct {
	NickName        string `json:"nickName"`
	Description     string `json:"description"`
	Avatar          string `json:"avatar"`
	PrivateGameData bool   `json:"privateGameData"`
}

type UpdateUserGameRequest struct {
	Status        int        `json:"status"`
	WishListMark  bool       `json:"wishListMark"`
	BlackListMark bool       `json:"blackListMark"`
	StartDate     *time.Time `json:"startDate"`
	FinishedDate  *time.Time `json:"finishedDate"`
}

type CreateUserGameRequest struct {
	GameErogsID   int        `json:"gameErogsId"`
	Status        int        `json:"status"`
	WishListMark  bool       `json:"wishListMark"`
	BlackListMark bool       `json:"blackListMark"`
	StartDate     *time.Time `json:"startDate"`
	FinishedDate  *time.Time `json:"finishedDate"`
}
