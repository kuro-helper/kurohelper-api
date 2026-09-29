package dto

import "time"

type EnsureGameRequest struct {
	ErogsID int `json:"erogsId"`
}

type UpdateGameRequest struct {
	ImageURL string `json:"imageUrl"`
}

type GameResponse struct {
	ID              int       `json:"id"`
	ErogsID         *int      `json:"erogsId"`
	ImageURL        string    `json:"imageUrl"`
	UpdatedUserName string    `json:"updatedUserName"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
