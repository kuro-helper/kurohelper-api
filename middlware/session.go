package middlware

import (
	"kurohelper-api/dto"
	"kurohelper-api/session"

	"github.com/gofiber/fiber/v3"
)

func SessionAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		if session.LoadUser(c) == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
				Message: "未登入",
				Data:    nil,
			})
		}

		return c.Next()
	}
}
