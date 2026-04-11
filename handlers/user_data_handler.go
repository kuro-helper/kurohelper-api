package handlers

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	kurohelperdb "github.com/peter910820/kurohelper-db/v2"
)

func GetUserHasPlayedHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userHasPlayed, err := kurohelperdb.SelectUserHasPlayed(id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "search successfully",
		"data":    userHasPlayed,
	})
}
