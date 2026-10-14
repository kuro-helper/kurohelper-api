package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// Statistics 統計資源
func StatisticsRouter(apiGroup fiber.Router) {
	group := apiGroup.Group("/statistics")

	group.Get("/user-game/finished", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetFinishedUserGamesHandler(c)
	})
}
