package middlware

import (
	"kurohelper-api/dto"
	"kurohelper-api/session"
	"log/slog"
	"strconv"

	"github.com/gofiber/fiber/v3"

	"kurohelperservice/db"
)

// RequireStaffRole 要求目前登入者為開發者或站主
func RequireStaffRole() fiber.Handler {
	return func(c fiber.Ctx) error {
		me := session.LoadUser(c)
		if me == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
				Message: "未登入",
				Data:    nil,
			})
		}

		user, err := db.GetUser(db.Dbs, strconv.Itoa(me.ID))
		if err != nil {
			slog.Error("RequireStaffRole GetUser", "err", err, "userId", me.ID)
			return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
				Message: "發生錯誤，請稍後再試",
				Data:    nil,
			})
		}

		if user.Role != db.UserRoleDeveloper && user.Role != db.UserRoleOwner {
			slog.Warn("RequireStaffRole forbidden", "userId", me.ID, "role", user.Role)
			return c.Status(fiber.StatusForbidden).JSON(dto.TResponse[any]{
				Message: "權限不足",
				Data:    nil,
			})
		}

		return c.Next()
	}
}
