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
}
