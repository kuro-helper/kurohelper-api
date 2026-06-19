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
	userDataGroup.Get("/", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUser(c)
	})

	userDataGroup.Get("/:id/game", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUserGameHandler(c)
	})

	userDataGroup.Get("/register-link", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetRegisterLinkHandler(c)
	})

	userDataGroup.Post("/register", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.RegisterUserHandler(c)
	})
}
