package handler

import (
	"kurohelper-api/dto"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"kurohelperservice/db"
)

func GetFinishedUserGamesHandler(c fiber.Ctx) error {
	games, err := db.GetUserGamesByStatus(db.Dbs, db.UserGameStatusFinished)
	if err != nil {
		slog.Error("GetUserGamesByStatus", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	privateUsers, err := db.GetUsersByPrivateGameData(db.Dbs, true)
	if err != nil {
		slog.Error("GetUsersByPrivateGameData", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	privateIDs := make(map[int]struct{}, len(privateUsers))
	for _, user := range privateUsers {
		privateIDs[user.ID] = struct{}{}
	}

	filtered := make([]db.UserGame, 0, len(games))
	for _, game := range games {
		if _, hidden := privateIDs[game.UserID]; hidden {
			continue
		}
		filtered = append(filtered, game)
	}

	slog.Info("GetFinishedUserGamesHandler success", "count", len(filtered))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.UserGameResponse]{
		Message: "ok",
		Data:    toUserGamesResponse(filtered),
	})
}
