package dto

import "time"

type CreateAnnouncementRequest struct {
	Category  string  `json:"category"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Icon      *string `json:"icon"`
	Thumbnail *string `json:"thumbnail"`
	Image     *string `json:"image"`
}

type AnnouncementResponse struct {
	ID        int       `json:"id"`
	Category  string    `json:"category"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Icon      *string   `json:"icon,omitempty"`      // MDI 名稱，如 mdi-bullhorn
	Thumbnail *string   `json:"thumbnail,omitempty"` // Discord側邊小圖 URL
	Image     *string   `json:"image,omitempty"`     // Discord底部大圖 URL
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
