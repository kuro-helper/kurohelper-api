package handler

import (
	"errors"
	"kurohelper-api/dto"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"kurohelperservice/db"
)

func GetUser(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	getData := func() ([]db.User, error) {
		if strings.TrimSpace(id) != "" {
			u, err := db.GetUserByDiscordID(db.Dbs, id)
			return []db.User{u}, err
		}
		return db.GetAllUsers(db.Dbs)
	}

	users, err := getData()
	if err != nil {
		slog.Error("GetUser", "err", err, "id", id)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "找不到該使用者",
				Data:    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	var userReturn []dto.User
	for _, u := range users {
		userReturn = append(userReturn, dto.User{
			ID:        u.DiscordID,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	slog.Info("GetUser success", "count", len(userReturn))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.User]{
		Message: "ok",
		Data:    userReturn,
	})
}

func GetUserHasPlayedHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userGames, err := db.GetUserGameByDiscordID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userHasPlayedResp := make([]dto.UserHasPlayed, 0, len(userGames))
	for _, game := range userGames {
		if game.Status != 1 {
			continue
		}

		dtoItem := dto.UserHasPlayed{
			UserID:      id,
			GameErogsID: game.GameErogsID,
			CompletedAt: game.FinishedDate,
			CreatedAt:   game.CreatedAt,
		}

		if game.GameErogs != nil {
			dtoItem.GameID = game.GameErogs.ID
			dtoItem.BrandID = game.GameErogs.BrandErogsID
			dtoItem.GameName = game.GameErogs.Name
			dtoItem.GameImage = game.GameErogs.Image

			if game.GameErogs.BrandErogs != nil {
				dtoItem.BrandName = game.GameErogs.BrandErogs.Name
				dtoItem.Disband = game.GameErogs.BrandErogs.Disband
				dtoItem.BrandGameCount = game.GameErogs.BrandErogs.GameCount
			}
		}

		userHasPlayedResp = append(userHasPlayedResp, dtoItem)
	}

	slog.Info("GetUserHasPlayedHandler success", "id", id, "count", len(userHasPlayedResp))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.UserHasPlayed]{
		Message: "search successfully",
		Data:    userHasPlayedResp,
	})
}

func GetUserInWishHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userGames, err := db.GetUserGameByDiscordID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserInWish", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userInWishResp := make([]dto.UserInWish, 0, len(userGames))
	for _, game := range userGames {
		if !game.WishListMark {
			continue
		}

		dtoItem := dto.UserInWish{
			UserID:      id,
			GameErogsID: game.GameErogsID,
			CreatedAt:   game.CreatedAt,
		}

		if game.GameErogs != nil {
			dtoItem.GameID = game.GameErogs.ID
			dtoItem.BrandID = game.GameErogs.BrandErogsID
			dtoItem.GameName = game.GameErogs.Name
			dtoItem.GameImage = game.GameErogs.Image

			if game.GameErogs.BrandErogs != nil {
				dtoItem.BrandName = game.GameErogs.BrandErogs.Name
				dtoItem.Disband = game.GameErogs.BrandErogs.Disband
				dtoItem.BrandGameCount = game.GameErogs.BrandErogs.GameCount
			}
		}

		userInWishResp = append(userInWishResp, dtoItem)
	}

	slog.Info("GetUserInWishHandler success", "id", id, "count", len(userInWishResp))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.UserInWish]{
		Message: "search successfully",
		Data:    userInWishResp,
	})
}

// 舊版API Handler
func GetUserHasPlayedLegacyHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userGames, err := db.GetUserGameByDiscordID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	userHasPlayed := make([]db.UserGame, 0, len(userGames))
	for _, game := range userGames {
		if game.Status == 1 {
			userHasPlayed = append(userHasPlayed, game)
		}
	}

	slog.Info("GetUserHasPlayedLegacyHandler success", "id", id, "count", len(userHasPlayed))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "search successfully",
		"data":    userHasPlayed,
	})
}
