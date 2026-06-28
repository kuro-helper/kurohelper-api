package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

func AuthRouter(apiGroup fiber.Router) {
	authGroup := apiGroup.Group("/auth")

	authGroup.Post("/login", middlware.TokenAuth(true), func(c fiber.Ctx) error {
		return handler.LoginHandler(c)
	})

	authGroup.Get("/register-link", middlware.TokenAuth(true), func(c fiber.Ctx) error {
		return handler.GetRegisterLinkHandler(c)
	})

	authGroup.Post("/register", middlware.TokenAuth(true), func(c fiber.Ctx) error {
		return handler.RegisterUserHandler(c)
	})

	authGroup.Get("/me", middlware.TokenAuth(true), func(c fiber.Ctx) error {
		return handler.MeHandler(c)
	})

	authGroup.Post("/logout", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.LogoutHandler(c)
	})
}
