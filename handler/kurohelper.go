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

func toGameResponse(item db.Game) dto.GameResponse {
	resp := dto.GameResponse{
		ID:        item.ID,
		ErogsID:   item.ErogsID,
		ImageURL:  item.ImageURL,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	if item.UpdatedUser > 0 {
		if user, err := db.GetUser(db.Dbs, strconv.Itoa(item.UpdatedUser)); err == nil {
			resp.UpdatedUserName = user.Name
		}
	}
	return resp
}

func EnsureGameHandler(c fiber.Ctx) error {
	me := session.LoadUser(c)
	if me == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
			Message: "未登入",
			Data:    nil,
		})
	}

	var req dto.EnsureGameRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("EnsureGameHandler bad request", "reason", "invalid body", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "請求格式錯誤",
			Data:    nil,
		})
	}
	if req.ErogsID <= 0 {
		slog.Warn("EnsureGameHandler bad request", "reason", "invalid erogsId", "erogsId", req.ErogsID)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "erogsId 必須是大於 0 的整數",
			Data:    nil,
		})
	}

	item, err := db.EnsureGameByErogsID(db.Dbs, req.ErogsID, me.ID)
	if err != nil {
		if errors.Is(err, db.ErrParameterNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
				Message: "erogsId 必須是大於 0 的整數",
				Data:    nil,
			})
		}
		slog.Error("EnsureGameByErogsID", "err", err, "erogsId", req.ErogsID, "userId", me.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("EnsureGameHandler success", "id", item.ID, "erogsId", req.ErogsID, "userId", me.ID)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.GameResponse]{
		Message: "ok",
		Data:    toGameResponse(item),
	})
}

func UpdateGameHandler(c fiber.Ctx) error {
	me := session.LoadUser(c)
	if me == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
			Message: "未登入",
			Data:    nil,
		})
	}

	idText := strings.TrimSpace(c.Params("id"))
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		slog.Warn("UpdateGameHandler bad request", "reason", "invalid id", "id", idText)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "id 必須是大於 0 的整數",
			Data:    nil,
		})
	}

	var req dto.UpdateGameRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("UpdateGameHandler bad request", "reason", "invalid body", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "請求格式錯誤",
			Data:    nil,
		})
	}
	imageURL := strings.TrimSpace(req.ImageURL)

	if err := db.UpdateGameImageURL(db.Dbs, id, imageURL, me.ID); err != nil {
		if errors.Is(err, db.ErrNoRowsAffected) || errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("UpdateGameHandler not found", "id", id, "userId", me.ID)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "找不到遊戲",
				Data:    nil,
			})
		}
		slog.Error("UpdateGameImageURL", "err", err, "id", id, "userId", me.ID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	item, err := db.GetGameByID(db.Dbs, id)
	if err != nil {
		slog.Error("GetGameByID", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("UpdateGameHandler success", "id", id, "userId", me.ID)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.GameResponse]{
		Message: "ok",
		Data:    toGameResponse(item),
	})
}
