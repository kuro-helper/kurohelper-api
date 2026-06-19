package middlware

import (
	"strings"

	kurohelperdb "kurohelperservice/db"

	"github.com/gofiber/fiber/v3"
)

var (
	VaildToken = make(map[string]kurohelperdb.WebAPIToken)
)

// 驗證有無合法Token
func TokenAuth(requirePrivileged bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Bearer格式檢查
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		token := parts[1]

		tokenData, ok := VaildToken[token]
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		if requirePrivileged && !tokenData.IsPrivileged {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}

		return c.Next()
	}
}
