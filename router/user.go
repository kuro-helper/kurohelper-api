package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// User資源獲取相關新版API
func UserRouter(apiGroup fiber.Router) {
	userDataGroup := apiGroup.Group("/user")

	// 獲取所有存在的User
	userDataGroup.Get("/", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetUser(c)
	})

	userDataGroup.Get("/:id/game", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetUserGameHandler(c)
	})

	userDataGroup.Post("/:id/game", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.CreateUserGameHandler(c)
	})

	userDataGroup.Put("/:id/game/:gameErogsId", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.UpdateUserGameHandler(c)
	})

	userDataGroup.Put("/:id", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.UpdateUserHandler(c)
	})
}
