package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// Kurohelper 自建資源
func KurohelperRouter(apiGroup fiber.Router) {
	group := apiGroup.Group("/kurohelper")

	group.Post("/game/ensure", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.EnsureGameHandler(c)
	})

	group.Put("/game/:id", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.UpdateGameHandler(c)
	})
}
