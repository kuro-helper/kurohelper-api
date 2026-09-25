package handler

import (
	"kurohelper-api/dto"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	"kurohelperservice/db"
)

func GetErogsBrand(c fiber.Ctx) error {
	brands, err := db.GetAllBrandErogs(db.Dbs)
	if err != nil {
		slog.Error("GetErogsBrand", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	out := make([]dto.BrandErogsResponse, 0, len(brands))
	for _, b := range brands {
		out = append(out, dto.BrandErogsResponse{
			ID:        b.ID,
			Name:      b.Name,
			Disband:   b.Disband,
			GameCount: b.GameCount,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	slog.Info("GetErogsBrand success", "count", len(out))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.BrandErogsResponse]{
		Message: "ok",
		Data:    out,
	})
}

func GetErogsGame(c fiber.Ctx) error {
	brandIDText := strings.TrimSpace(c.Query("brandid"))
	brandID, err := strconv.Atoi(brandIDText)
	if err != nil || brandID <= 0 {
		slog.Warn("GetErogsGame bad request", "reason", "invalid brandid", "brandid", brandIDText)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "brandid 必須是大於 0 的整數",
			Data:    nil,
		})
	}

	games, err := db.GetGameErogsByBrandID(db.Dbs, brandID)
	if err != nil {
		slog.Error("GetErogsGame", "err", err, "brandid", brandID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	out := make([]dto.GameErogsResponse, 0, len(games))
	for _, g := range games {
		out = append(out, dto.GameErogsResponse{
			ID:           g.ID,
			BrandErogsID: g.BrandErogsID,
			Name:         g.Name,
			Image:        g.Image,
			Category:     g.Category,
			CreatedAt:    g.CreatedAt,
			UpdatedAt:    g.UpdatedAt,
		})
	}

	slog.Info("GetErogsGame success", "brandid", brandID, "count", len(out))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.GameErogsResponse]{
		Message: "ok",
		Data:    out,
	})
}
