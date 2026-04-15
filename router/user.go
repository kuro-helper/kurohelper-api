package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

func UserRouter(apiGroup fiber.Router) {
	userDataGroup := apiGroup.Group("/user")

	// 獲取所有存在的User
	userDataGroup.Get("/", middlware.TokenAuth(), func(c fiber.Ctx) error {
		return handler.GetUser(c)
	})
}
