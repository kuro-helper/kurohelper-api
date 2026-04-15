package handler

import (
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	// 寫到快取
	middlware.VaildToken[id.String()] = struct{}{}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "token generated successfully",
		"token":   id.String(),
	})
}
