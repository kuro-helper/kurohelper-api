package middlware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// SessionUserIDKey is the session field set on successful user login.
const SessionUserIDKey = "userId"

func UserAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		sess := session.FromContext(c)
		if sess == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "session not available",
			})
		}

		if sess.Get(SessionUserIDKey) == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "not logged in",
			})
		}

		return c.Next()
	}
}
