package dto

import "time"

type CreateArticleRequest struct {
	GameErogsID int    `json:"gameErogsId"`
	Title       string `json:"title"`
	Content     string `json:"content"`
}

type ArticleUserResponse struct {
	ID       int    `json:"id"`
	NickName string `json:"nickName"`
	Avatar   string `json:"avatar"`
}

type ArticleResponse struct {
	UserID      int                  `json:"userId"`
	GameErogsID int                  `json:"gameErogsId"`
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	User        *ArticleUserResponse `json:"user,omitempty"`
	GameErogs   *GameErogsResponse   `json:"gameErogs,omitempty"`
}
