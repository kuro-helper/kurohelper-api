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
			u, err := db.GetUser(db.Dbs, id)
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
			ID:        u.ID,
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

	userHasPlayed, err := db.GetUserHasPlayedByID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userHasPlayedResp := make([]dto.UserHasPlayed, 0, len(userHasPlayed))
	for _, hasPlayed := range userHasPlayed {
		dtoItem := dto.UserHasPlayed{
			UserID:      hasPlayed.UserID,
			GameErogsID: hasPlayed.GameErogsID,
			CompletedAt: hasPlayed.CompletedAt,
			CreatedAt:   hasPlayed.CreatedAt,
		}

		if hasPlayed.GameErogs != nil {
			dtoItem.GameID = hasPlayed.GameErogs.ID
			dtoItem.BrandID = hasPlayed.GameErogs.BrandErogsID
			dtoItem.GameName = hasPlayed.GameErogs.Name
			dtoItem.GameImage = hasPlayed.GameErogs.Image

			if hasPlayed.GameErogs.BrandErogs != nil {
				dtoItem.BrandName = hasPlayed.GameErogs.BrandErogs.Name
				dtoItem.Disband = hasPlayed.GameErogs.BrandErogs.Disband
				dtoItem.BrandGameCount = hasPlayed.GameErogs.BrandErogs.GameCount
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

	userInWish, err := db.GetUserInWishByID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserInWish", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userInWishResp := make([]dto.UserInWish, 0, len(userInWish))
	for _, inWish := range userInWish {
		dtoItem := dto.UserInWish{
			UserID:      inWish.UserID,
			GameErogsID: inWish.GameErogsID,
			CreatedAt:   inWish.CreatedAt,
		}

		if inWish.GameErogs != nil {
			dtoItem.GameID = inWish.GameErogs.ID
			dtoItem.BrandID = inWish.GameErogs.BrandErogsID
			dtoItem.GameName = inWish.GameErogs.Name
			dtoItem.GameImage = inWish.GameErogs.Image

			if inWish.GameErogs.BrandErogs != nil {
				dtoItem.BrandName = inWish.GameErogs.BrandErogs.Name
				dtoItem.Disband = inWish.GameErogs.BrandErogs.Disband
				dtoItem.BrandGameCount = inWish.GameErogs.BrandErogs.GameCount
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

	userHasPlayed, err := db.GetUserHasPlayedByID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	slog.Info("GetUserHasPlayedLegacyHandler success", "id", id, "count", len(userHasPlayed))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "search successfully",
		"data":    userHasPlayed,
	})
}
