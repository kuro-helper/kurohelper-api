package middlware

import (
	"encoding/gob"
	"kurohelper-api/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

// SessionUserKey is the session field for the logged-in user snapshot.
const SessionUserKey = "user"

func init() {
	gob.Register(dto.SessionUser{})
}

// 驗證有無登入Session狀態
func SessionAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		if _, ok := sessionUserFromContext(c); !ok {
			sess := session.FromContext(c)
			if sess == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "session not available",
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "not logged in",
			})
		}

		return c.Next()
	}
}

func sessionUserFromContext(c fiber.Ctx) (dto.SessionUser, bool) {
	sess := session.FromContext(c)
	if sess == nil {
		return dto.SessionUser{}, false
	}

	raw := sess.Get(SessionUserKey)
	if raw == nil {
		return dto.SessionUser{}, false
	}

	user, ok := raw.(dto.SessionUser)
	return user, ok
}
