package handler

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"

	kurohelperdb "kurohelperservice/db"
)

func GetUserHasPlayedHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userHasPlayed, err := kurohelperdb.GetUserHasPlayedByID(kurohelperdb.Dbs, id)
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
