package routes

import (
	"kurohelper-api/handlers"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

func UserDataRouter(apiGroup fiber.Router) {
	userDataGroup := apiGroup.Group("/userdata")

	// 獲取指定使用者全部的遊玩資料
	userDataGroup.Get("/", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handlers.GetUserHasPlayedHandler(c)
	})
}
