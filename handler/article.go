package handler

import (
	"errors"
	"kurohelper-api/dto"
	"kurohelper-api/session"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"kurohelperservice/db"
)

func toArticleResponse(article db.UserGameArticle) dto.ArticleResponse {
	response := dto.ArticleResponse{
		UserID:      article.UserID,
		GameErogsID: article.GameErogsID,
		Title:       article.Title,
		Content:     article.Content,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
	if article.User != nil {
		response.User = &dto.ArticleUserResponse{
			ID:       article.User.ID,
			NickName: article.User.Name,
			Avatar:   article.User.Avatar,
		}
	}
	if game := article.GameErogs; game != nil {
		response.GameErogs = &dto.GameErogsResponse{
			ID:           game.ID,
			BrandErogsID: game.BrandErogsID,
			Name:         game.Name,
			Image:        game.Image,
			Category:     game.Category,
			CreatedAt:    game.CreatedAt,
			UpdatedAt:    game.UpdatedAt,
		}
	}
	return response
}

func GetArticlesHandler(c fiber.Ctx) error {
	userIDText := strings.TrimSpace(c.Query("userId"))
	gameErogsIDText := strings.TrimSpace(c.Query("gameErogsId"))
	var userID, gameErogsID int
	var err error
	if userIDText != "" {
		userID, err = strconv.Atoi(userIDText)
		if err != nil || userID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{Message: "userId 須為正整數", Data: nil})
		}
	}
	if gameErogsIDText != "" {
		gameErogsID, err = strconv.Atoi(gameErogsIDText)
		if err != nil || gameErogsID <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{Message: "gameErogsId 須為正整數", Data: nil})
		}
	}

	var articles []db.UserGameArticle
	switch {
	case userID > 0 && gameErogsID > 0:
		articles, err = db.GetUserGameArticlesByUserID(db.Dbs.Where("user_game_articles.game_erogs_id = ?", gameErogsID), userID)
	case userID > 0:
		articles, err = db.GetUserGameArticlesByUserID(db.Dbs, userID)
	case gameErogsID > 0:
		articles, err = db.GetUserGameArticlesByGameErogsID(db.Dbs, gameErogsID)
	default:
		articles, err = db.GetAllUserGameArticles(db.Dbs)
	}
	if err != nil {
		slog.Error("GetArticles", "err", err, "userId", userID, "gameErogsId", gameErogsID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "載入文章失敗，請稍後再試",
			Data:    nil,
		})
	}

	viewerID := sessionUserID(c)
	response := make([]dto.ArticleResponse, 0, len(articles))
	for _, article := range articles {
		if article.User == nil || (article.User.PrivateGameData && article.UserID != viewerID) {
			continue
		}
		response = append(response, toArticleResponse(article))
	}
	return c.JSON(dto.TResponse[[]dto.ArticleResponse]{Message: "ok", Data: response})
}

func CreateArticleHandler(c fiber.Ctx) error {
	user := session.LoadUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{Message: "未登入", Data: nil})
	}
	var request dto.CreateArticleRequest
	if err := c.Bind().Body(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{Message: "請求格式錯誤", Data: nil})
	}
	request.Title = strings.TrimSpace(request.Title)
	if request.GameErogsID <= 0 || request.Title == "" || strings.TrimSpace(request.Content) == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "gameErogsId 須為正整數，title 與 content 不可為空",
			Data:    nil,
		})
	}
	if err := db.CreateUserGameArticle(db.Dbs, user.ID, request.GameErogsID, request.Title, request.Content); err != nil {
		if errors.Is(err, db.ErrUniqueViolation) {
			return c.Status(fiber.StatusConflict).JSON(dto.TResponse[any]{Message: "你已為此遊戲建立文章", Data: nil})
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{Message: "找不到該遊戲", Data: nil})
		}
		slog.Error("CreateUserGameArticle", "err", err, "userId", user.ID, "gameErogsId", request.GameErogsID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{Message: "建立文章失敗，請稍後再試", Data: nil})
	}
	return c.Status(fiber.StatusCreated).JSON(dto.TResponse[any]{Message: "ok", Data: nil})
}
