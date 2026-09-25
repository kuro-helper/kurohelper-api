package dto

import "time"

type BrandErogsResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Disband   bool      `json:"disband"`
	GameCount int       `json:"gameCount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type GameErogsResponse struct {
	ID           int       `json:"id"`
	BrandErogsID int       `json:"brandErogsId"`
	Name         string    `json:"name"`
	Image        string    `json:"image"`
	Category     string    `json:"category"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
