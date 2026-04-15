package handler

import (
	"kurohelper-api/dto"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"kurohelperservice/db"
)

func GetUser(c fiber.Ctx) error {
	user, err := db.GetAllUsers(db.Dbs)
	if err != nil {
		slog.Error("GetUser", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "發生錯誤，請稍後再試",
		})
	}

	var userReturn []dto.User
	for _, u := range user {
		userReturn = append(userReturn, GetUserMapper(&u))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": userReturn,
	})

}

func GetUserMapper(user *db.User) dto.User {
	return dto.User{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
