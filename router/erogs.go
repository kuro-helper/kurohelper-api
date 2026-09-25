package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// Erogs 相關資源
func ErogsRouter(apiGroup fiber.Router) {
	group := apiGroup.Group("/erogs")

	group.Get("/brand/", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetErogsBrand(c)
	})

	group.Get("/game/", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetErogsGame(c)
	})
}
