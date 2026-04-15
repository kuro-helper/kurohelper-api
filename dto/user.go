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
