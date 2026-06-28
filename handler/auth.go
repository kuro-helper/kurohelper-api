package handler

import (
	"errors"
	"kurohelper-api/dto"
	"kurohelper-api/session"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
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

	sessionUser := session.NewUser(user, auth.Username)
	if !session.SetUser(c, sessionUser) {
		slog.Error("LoginHandler session not available")
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "登入失敗，請稍後再試",
			Data:    nil,
		})
	}

	meUser := toAuthUser(user, auth.Username)
	slog.Info("LoginHandler success", "userId", meUser.ID, "userName", auth.Username)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.LoginResponse]{
		Message: "ok",
		Data: dto.LoginResponse{
			User: meUser,
		},
	})
}

func MeHandler(c fiber.Ctx) error {
	user := session.LoadUser(c)
	if user == nil {
		return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.MeResponse]{
			Message: "ok",
			Data:    dto.MeResponse{User: nil},
		})
	}

	meUser := authUserFromSession(*user)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.MeResponse]{
		Message: "ok",
		Data: dto.MeResponse{
			User: &meUser,
		},
	})
}

func authUserFromSession(user session.User) dto.AuthUser {
	return dto.AuthUser{
		ID:          user.ID,
		UserName:    user.UserName,
		NickName:    user.NickName,
		DiscordID:   user.DiscordID,
		Avatar:      user.Avatar,
		Description: user.Description,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}

func toAuthUser(u db.User, userName string) dto.AuthUser {
	return dto.AuthUser{
		ID:          u.ID,
		UserName:    userName,
		NickName:    u.Name,
		DiscordID:   discordIDOrEmpty(u.DiscordID),
		Avatar:      u.Avatar,
		Description: u.Description,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
