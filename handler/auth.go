package handler

import (
	"errors"
	"kurohelper-api/dto"
	"kurohelper-api/middlware"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kurohelperservice/db"
)

func LoginHandler(c fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("LoginHandler bad request", "reason", "bind body", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "請求格式錯誤",
			Data:    nil,
		})
	}

	req.UserName = strings.TrimSpace(req.UserName)
	req.Password = strings.TrimSpace(req.Password)

	if req.UserName == "" || req.Password == "" {
		slog.Warn("LoginHandler bad request", "reason", "empty fields")
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "userName、password 不可為空",
			Data:    nil,
		})
	}

	auth, err := db.GetUserAuthByUsername(db.Dbs, req.UserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("LoginHandler unauthorized", "userName", req.UserName)
			return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
				Message: "帳號或密碼錯誤",
				Data:    nil,
			})
		}
		slog.Error("GetUserAuthByUsername", "err", err, "userName", req.UserName)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "登入失敗，請稍後再試",
			Data:    nil,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(req.Password)); err != nil {
		slog.Warn("LoginHandler unauthorized", "userName", req.UserName, "reason", "password mismatch")
		return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
			Message: "帳號或密碼錯誤",
			Data:    nil,
		})
	}

	sess := session.FromContext(c)
	if sess == nil {
		slog.Error("LoginHandler session not available")
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "登入失敗，請稍後再試",
			Data:    nil,
		})
	}

	user, err := db.GetUser(db.Dbs, strconv.Itoa(auth.UserID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("LoginHandler user not found", "userId", auth.UserID)
			return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
				Message: "帳號或密碼錯誤",
				Data:    nil,
			})
		}
		slog.Error("GetUser", "err", err, "userId", auth.UserID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "登入失敗，請稍後再試",
			Data:    nil,
		})
	}

	sessionUser := toSessionUser(user)
	sess.Set(middlware.SessionUserKey, sessionUser)

	slog.Info("LoginHandler success", "userId", sessionUser.ID, "userName", auth.Username)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.LoginResponse]{
		Message: "ok",
		Data: dto.LoginResponse{
			User: sessionUser,
		},
	})
}

func toSessionUser(u db.User) dto.SessionUser {
	return dto.SessionUser{
		ID:          u.ID,
		Name:        u.Name,
		DiscordID:   discordIDOrEmpty(u.DiscordID),
		Avatar:      u.Avatar,
		Description: u.Description,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
