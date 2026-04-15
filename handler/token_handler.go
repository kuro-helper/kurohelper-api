package handler

import (
	"kurohelper-api/dto"
	"kurohelper-api/middlware"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	kurohelperdb "kurohelperservice/db"
)

// internal environment tokens generate handler
func TokensGenerateHandler(c fiber.Ctx) error {
	id := uuid.New()

	// 寫到db(目前過期時間預設都是無限，因為只有內部使用)
	err := kurohelperdb.CreateWebAPIToken(kurohelperdb.Dbs, id.String(), 0)
	if err != nil {
		slog.Error("CreateWebAPIToken", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	// 寫到快取
	middlware.VaildToken[id.String()] = struct{}{}

	slog.Info("TokensGenerateHandler success", "token", id.String())
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[string]{
		Message: "token generated successfully",
		Data:    id.String(),
	})
}
