package handler

import (
	"errors"
	"kurohelper-api/cache"
	"kurohelper-api/dto"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"

	kurohelperservice "kurohelperservice"
	"kurohelperservice/db"
	"kurohelperservice/provider/erogs"
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
		image, imageFromGame := resolveErogsDisplayImage(g.Image, g.GameImageURL)
		out = append(out, dto.GameErogsResponse{
			ID:            g.ID,
			BrandErogsID:  g.BrandErogsID,
			Name:          g.Name,
			Image:         image,
			ImageFromGame: imageFromGame,
			Category:      g.Category,
			CreatedAt:     g.CreatedAt,
			UpdatedAt:     g.UpdatedAt,
		})
	}

	slog.Info("GetErogsGame success", "brandid", brandID, "count", len(out))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.GameErogsResponse]{
		Message: "ok",
		Data:    out,
	})
}

func GetErogsGameByID(c fiber.Ctx) error {
	idText := strings.TrimSpace(c.Params("id"))
	id, err := strconv.Atoi(idText)
	if err != nil || id <= 0 {
		slog.Warn("GetErogsGameByID bad request", "reason", "invalid id", "id", idText)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "id 必須是大於 0 的整數",
			Data:    nil,
		})
	}

	cacheKey := strconv.Itoa(id)
	game, err := cache.ErogsGameStore.Get(cacheKey)
	if err != nil {
		if !errors.Is(err, kurohelperservice.ErrCacheLost) {
			slog.Error("GetErogsGameByID cache", "err", err, "id", id)
			return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
				Message: "發生錯誤，請稍後再試",
				Data:    nil,
			})
		}

		game, err = erogs.SearchGameByID(id)
		if err != nil {
			if errors.Is(err, kurohelperservice.ErrSearchNoContent) {
				slog.Warn("GetErogsGameByID not found", "id", id)
				return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
					Message: "找不到遊戲",
					Data:    nil,
				})
			}
			if errors.Is(err, kurohelperservice.ErrRateLimit) {
				slog.Warn("GetErogsGameByID rate limited", "id", id)
				return c.Status(fiber.StatusTooManyRequests).JSON(dto.TResponse[any]{
					Message: "批評空間查詢過於頻繁，請稍後再試",
					Data:    nil,
				})
			}
			slog.Error("GetErogsGameByID", "err", err, "id", id)
			return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
				Message: "發生錯誤，請稍後再試",
				Data:    nil,
			})
		}
		if game == nil || game.ID == 0 {
			slog.Warn("GetErogsGameByID empty result", "id", id)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "找不到遊戲",
				Data:    nil,
			})
		}
		cache.ErogsGameStore.Set(cacheKey, game)
	} else {
		slog.Info("GetErogsGameByID cache hit", "id", id)
	}

	creators := make([]dto.ErogsOfficialGameCreatorResponse, 0, len(game.CreatorShubetu))
	for _, item := range game.CreatorShubetu {
		creators = append(creators, dto.ErogsOfficialGameCreatorResponse{
			ShubetuType:       item.ShubetuType,
			CreatorName:       item.CreatorName,
			ShubetuDetailType: item.ShubetuDetailType,
			ShubetuDetailName: item.ShubetuDetailName,
		})
	}

	slog.Info("GetErogsGameByID success", "id", id, "name", game.Gamename)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.ErogsOfficialGameResponse]{
		Message: "ok",
		Data: dto.ErogsOfficialGameResponse{
			ID:                               game.ID,
			BrandID:                          game.BrandID,
			BrandName:                        game.BrandName,
			Name:                             game.Gamename,
			SellDay:                          game.SellDay,
			Model:                            game.Model,
			DMM:                              game.DMM,
			Median:                           game.Median,
			TokutenCount:                     game.TokutenCount,
			TotalPlayTimeMedian:              game.TotalPlayTimeMedian,
			TimeBeforeUnderstandingFunMedian: game.TimeBeforeUnderstandingFunMedian,
			Okazu:                            game.Okazu,
			Erogame:                          game.Erogame,
			Genre:                            game.Genre,
			BannerURL:                        game.BannerUrl,
			SteamID:                          game.SteamId,
			VndbID:                           game.VndbId,
			Shoukai:                          game.Shoukai,
			Junni:                            game.Junni,
			Creators:                         creators,
		},
	})
}
