package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// 新版API
func UserRouter(apiGroup fiber.Router) {
	userDataGroup := apiGroup.Group("/user")

	// 獲取所有存在的User
	userDataGroup.Get("/", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUser(c)
	})

	userDataGroup.Get("/:id/game", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUserGameHandler(c)
	})

	// 獲取指定使用者全部的遊玩資料
	userDataGroup.Get("/played", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUserHasPlayedHandler(c)
	})

	// 獲取指定使用者全部的願望清單資料
	userDataGroup.Get("/wish", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUserInWishHandler(c)
	})

	userDataGroup.Get("/register-link", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetRegisterLinkHandler(c)
	})

	userDataGroup.Post("/register", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.RegisterUserHandler(c)
	})
}

// 舊版API
func UserLegacyRouter(apiGroup fiber.Router) {
	userDataGroup := apiGroup.Group("/userdata")

	// 獲取指定使用者全部的遊玩資料
	userDataGroup.Get("/", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUserHasPlayedLegacyHandler(c)
	})
}
