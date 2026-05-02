package handler

import (
	"errors"
	"kurohelper-api/dto"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"kurohelperservice/db"
)

func discordIDOrEmpty(discordID *string) string {
	if discordID == nil {
		return ""
	}
	return *discordID
}

func GetRegisterLinkHandler(c fiber.Ctx) error {
	registerID := strings.TrimSpace(c.Query("register_id"))
	if registerID == "" {
		slog.Warn("GetRegisterLinkHandler bad request", "reason", "empty register_id")
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "register_id 不可為空",
			Data:    nil,
		})
	}

	cacheData, err := db.GetRegisterCacheByID(db.Dbs, registerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("GetRegisterLinkHandler not found", "register_id", registerID)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "註冊連結不存在或已過期",
				Data:    nil,
			})
		}
		slog.Error("GetRegisterCacheByID", "err", err, "register_id", registerID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("GetRegisterLinkHandler success", "register_id", registerID)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.RegisterLookupResponse]{
		Message: "ok",
		Data: dto.RegisterLookupResponse{
			DiscordID: cacheData.DiscordID,
		},
	})
}

func RegisterUserHandler(c fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("RegisterUserHandler bad request", "reason", "bind body", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "請求格式錯誤",
			Data:    nil,
		})
	}

	req.RegisterID = strings.TrimSpace(req.RegisterID)
	req.UserName = strings.TrimSpace(req.UserName)
	req.Password = strings.TrimSpace(req.Password)

	if req.RegisterID == "" || req.UserName == "" || req.Password == "" {
		slog.Warn("RegisterUserHandler bad request", "reason", "empty fields",
			"register_id_present", req.RegisterID != "",
			"user_name_present", req.UserName != "",
			"password_present", req.Password != "",
		)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "register_id、user_name、password 不可為空",
			Data:    nil,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("GenerateFromPassword", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "密碼加密失敗，請稍後再試",
			Data:    nil,
		})
	}

	discordID := ""
	err = db.Dbs.Transaction(func(tx *gorm.DB) error {
		cacheData, err := db.GetRegisterCacheByID(tx, req.RegisterID)
		if err != nil {
			return err
		}
		discordID = cacheData.DiscordID

		user, err := db.EnsureDiscordUser(tx, cacheData.DiscordID, req.UserName)
		if err != nil {
			return err
		}

		// 確保username不重複
		existUserAuth, err := db.GetUserAuthByUsername(tx, req.UserName)
		if err == nil && existUserAuth.UserID != user.ID {
			return db.ErrUniqueViolation
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		// 確保不重複註冊
		_, err = db.GetUserAuthByUserID(tx, user.ID)
		if err == nil {
			return db.ErrUniqueViolation
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := db.CreateUserAuth(tx, user.ID, req.UserName, string(hashedPassword)); err != nil {
			return err
		}

		if err := db.DeleteRegisterCacheByID(tx, req.RegisterID); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("RegisterUserHandler not found", "register_id", req.RegisterID)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "註冊連結不存在或已過期",
				Data:    nil,
			})
		}
		if errors.Is(err, db.ErrUniqueViolation) {
			slog.Warn("RegisterUserHandler conflict", "register_id", req.RegisterID, "discord_id", discordID, "user_name", req.UserName)
			return c.Status(fiber.StatusConflict).JSON(dto.TResponse[any]{
				Message: "user_name 已存在或該 Discord 帳號已註冊",
				Data:    nil,
			})
		}
		slog.Error("RegisterUserHandler", "err", err, "register_id", req.RegisterID, "discord_id", discordID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "註冊失敗，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("RegisterUserHandler success", "discord_id", discordID, "user_name", req.UserName, "register_id", req.RegisterID)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.RegisterResponse]{
		Message: "register successfully",
		Data: dto.RegisterResponse{
			DiscordID: discordID,
			UserName:  req.UserName,
		},
	})
}

func GetUser(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	getData := func() ([]db.User, error) {
		if strings.TrimSpace(id) != "" {
			u, err := db.GetUserByDiscordID(db.Dbs, id)
			return []db.User{u}, err
		}
		return db.GetAllUsers(db.Dbs)
	}

	users, err := getData()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("GetUser not found", "id", id)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "找不到該使用者",
				Data:    nil,
			})
		}
		slog.Error("GetUser", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	var userReturn []dto.UserResponse
	for _, u := range users {
		userReturn = append(userReturn, dto.UserResponse{
			ID:        discordIDOrEmpty(u.DiscordID),
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	slog.Info("GetUser success", "id", id, "count", len(userReturn))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.UserResponse]{
		Message: "ok",
		Data:    userReturn,
	})
}

func GetUserHasPlayedHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userGames, err := db.GetUserGameByDiscordID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userHasPlayedResp := make([]dto.UserHasPlayedResponse, 0, len(userGames))
	for _, game := range userGames {
		if game.Status != 1 {
			continue
		}

		dtoItem := dto.UserHasPlayedResponse{
			UserID:      id,
			GameErogsID: game.GameErogsID,
			CompletedAt: game.FinishedDate,
			CreatedAt:   game.CreatedAt,
		}

		if game.GameErogs != nil {
			dtoItem.GameID = game.GameErogs.ID
			dtoItem.BrandID = game.GameErogs.BrandErogsID
			dtoItem.GameName = game.GameErogs.Name
			dtoItem.GameImage = game.GameErogs.Image

			if game.GameErogs.BrandErogs != nil {
				dtoItem.BrandName = game.GameErogs.BrandErogs.Name
				dtoItem.Disband = game.GameErogs.BrandErogs.Disband
				dtoItem.BrandGameCount = game.GameErogs.BrandErogs.GameCount
			}
		}

		userHasPlayedResp = append(userHasPlayedResp, dtoItem)
	}

	slog.Info("GetUserHasPlayedHandler success", "id", id, "count", len(userHasPlayedResp))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.UserHasPlayedResponse]{
		Message: "search successfully",
		Data:    userHasPlayedResp,
	})
}

func GetUserInWishHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userGames, err := db.GetUserGameByDiscordID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserInWish", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: err.Error(),
			Data:    nil,
		})
	}

	userInWishResp := make([]dto.UserInWishResponse, 0, len(userGames))
	for _, game := range userGames {
		if !game.WishListMark {
			continue
		}

		dtoItem := dto.UserInWishResponse{
			UserID:      id,
			GameErogsID: game.GameErogsID,
			CreatedAt:   game.CreatedAt,
		}

		if game.GameErogs != nil {
			dtoItem.GameID = game.GameErogs.ID
			dtoItem.BrandID = game.GameErogs.BrandErogsID
			dtoItem.GameName = game.GameErogs.Name
			dtoItem.GameImage = game.GameErogs.Image

			if game.GameErogs.BrandErogs != nil {
				dtoItem.BrandName = game.GameErogs.BrandErogs.Name
				dtoItem.Disband = game.GameErogs.BrandErogs.Disband
				dtoItem.BrandGameCount = game.GameErogs.BrandErogs.GameCount
			}
		}

		userInWishResp = append(userInWishResp, dtoItem)
	}

	slog.Info("GetUserInWishHandler success", "id", id, "count", len(userInWishResp))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.UserInWishResponse]{
		Message: "search successfully",
		Data:    userInWishResp,
	})
}

// 舊版API Handler
func GetUserHasPlayedLegacyHandler(c fiber.Ctx) error {
	// URL decoding
	id := c.Query("id")

	userGames, err := db.GetUserGameByDiscordID(db.Dbs, id)
	if err != nil {
		slog.Error("SelectUserHasPlayed", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	userHasPlayed := make([]db.UserGame, 0, len(userGames))
	for _, game := range userGames {
		if game.Status == 1 {
			userHasPlayed = append(userHasPlayed, game)
		}
	}

	slog.Info("GetUserHasPlayedLegacyHandler success", "id", id, "count", len(userHasPlayed))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "search successfully",
		"data":    userHasPlayed,
	})
}
