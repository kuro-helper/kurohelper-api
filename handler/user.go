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

func discordIDOrEmpty(discordID *string) string {
	if discordID == nil {
		return ""
	}
	return *discordID
}

func GetRegisterLinkHandler(c fiber.Ctx) error {
	registerID := strings.TrimSpace(c.Query("registerId"))
	if registerID == "" {
		slog.Warn("GetRegisterLinkHandler bad request", "reason", "empty registerId")
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "registerId 不可為空",
			Data:    nil,
		})
	}

	cacheData, err := db.GetRegisterCacheByID(db.Dbs, registerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("GetRegisterLinkHandler not found", "registerId", registerID)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "註冊連結不存在或已過期",
				Data:    nil,
			})
		}
		slog.Error("GetRegisterCacheByID", "err", err, "registerId", registerID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("GetRegisterLinkHandler success", "registerId", registerID)
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
			"registerIdPresent", req.RegisterID != "",
			"userNamePresent", req.UserName != "",
			"passwordPresent", req.Password != "",
		)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "registerId、userName、password 不可為空",
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

		var user db.User
		existing, err := db.GetUserByDiscordID(tx, cacheData.DiscordID)
		switch {
		case err == nil:
			user = existing
		case errors.Is(err, gorm.ErrRecordNotFound):
			created, err := db.EnsureDiscordUser(tx, cacheData.DiscordID, req.UserName)
			if err != nil {
				return err
			}
			user = *created
		default:
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
			slog.Warn("RegisterUserHandler not found", "registerId", req.RegisterID)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "註冊連結不存在或已過期",
				Data:    nil,
			})
		}
		if errors.Is(err, db.ErrUniqueViolation) {
			slog.Warn("RegisterUserHandler conflict", "registerId", req.RegisterID, "discordId", discordID, "userName", req.UserName)
			return c.Status(fiber.StatusConflict).JSON(dto.TResponse[any]{
				Message: "userName 已存在或該 Discord 帳號已註冊",
				Data:    nil,
			})
		}
		slog.Error("RegisterUserHandler", "err", err, "registerId", req.RegisterID, "discordId", discordID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "註冊失敗，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("RegisterUserHandler success", "discordId", discordID, "userName", req.UserName, "registerId", req.RegisterID)
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
			ID:        u.ID,
			NickName:  u.Name,
			DiscordID: discordIDOrEmpty(u.DiscordID),
			Role:      u.Role,
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

func profileAuthUserName(userID int) string {
	auth, err := db.GetUserAuthByUserID(db.Dbs, userID)
	if err != nil {
		return ""
	}
	return auth.Username
}

func toUserProfileResponse(u db.User) dto.UserProfileResponse {
	return dto.UserProfileResponse{
		ID:          u.ID,
		UserName:    profileAuthUserName(u.ID),
		NickName:    u.Name,
		DiscordID:   discordIDOrEmpty(u.DiscordID),
		Avatar:      u.Avatar,
		Description: u.Description,
		Role:        u.Role,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func userGameStatusText(status int) string {
	if status == 1 {
		return "finished"
	}
	return strconv.Itoa(status)
}

func toUserGameResponse(game db.UserGame) dto.UserGameResponse {
	item := dto.UserGameResponse{
		UserID:        game.UserID,
		GameErogsID:   game.GameErogsID,
		Status:        userGameStatusText(game.Status),
		WishListMark:  game.WishListMark,
		BlackListMark: game.BlackListMark,
		StartDate:     game.StartDate,
		FinishedDate:  game.FinishedDate,
		CreatedAt:     game.CreatedAt,
		UpdatedAt:     game.UpdatedAt,
	}

	if game.GameErogs == nil {
		return item
	}

	gameErogs := dto.UserGameErogsResponse{
		ID:           game.GameErogs.ID,
		BrandErogsID: game.GameErogs.BrandErogsID,
		Name:         game.GameErogs.Name,
		Image:        game.GameErogs.Image,
		CreatedAt:    game.GameErogs.CreatedAt,
		UpdatedAt:    game.GameErogs.UpdatedAt,
	}

	if game.GameErogs.BrandErogs != nil {
		brand := game.GameErogs.BrandErogs
		gameErogs.BrandErogs = &dto.UserGameBrandErogsResponse{
			ID:        brand.ID,
			Name:      brand.Name,
			Disband:   brand.Disband,
			GameCount: brand.GameCount,
			CreatedAt: brand.CreatedAt,
			UpdatedAt: brand.UpdatedAt,
		}
	}

	item.GameErogs = &gameErogs
	return item
}

func toUserGamesResponse(games []db.UserGame) []dto.UserGameResponse {
	resp := make([]dto.UserGameResponse, 0, len(games))
	for _, game := range games {
		resp = append(resp, toUserGameResponse(game))
	}
	return resp
}

func toGetUserGameResponse(u db.User, games []db.UserGame) dto.GetUserGameResponse {
	return dto.GetUserGameResponse{
		User:  toUserProfileResponse(u),
		Games: toUserGamesResponse(games),
	}
}

func GetUserGameHandler(c fiber.Ctx) error {
	id := strings.TrimSpace(c.Params("id"))
	userID, err := strconv.Atoi(id)
	if err != nil || userID <= 0 {
		slog.Warn("GetUserGameHandler bad request", "reason", "invalid id", "id", id)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "id 格式錯誤",
			Data:    nil,
		})
	}

	user, err := db.GetUser(db.Dbs, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("GetUserGameHandler not found", "id", id)
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

	userGames, err := db.GetUserGameByUserID(db.Dbs, userID)
	if err != nil {
		slog.Error("GetUserGameByUserID", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("GetUserGameHandler success", "id", id, "count", len(userGames))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.GetUserGameResponse]{
		Message: "ok",
		Data:    toGetUserGameResponse(user, userGames),
	})
}

func UpdateUserHandler(c fiber.Ctx) error {
	me := session.LoadUser(c)
	if me == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
			Message: "未登入",
			Data:    nil,
		})
	}

	id := strings.TrimSpace(c.Params("id"))
	userID, err := strconv.Atoi(id)
	if err != nil || userID <= 0 {
		slog.Warn("UpdateUserHandler bad request", "reason", "invalid id", "id", id)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "id 格式錯誤",
			Data:    nil,
		})
	}

	if me.ID != userID {
		slog.Warn("UpdateUserHandler forbidden", "sessionUserId", me.ID, "targetId", userID)
		return c.Status(fiber.StatusForbidden).JSON(dto.TResponse[any]{
			Message: "無法修改他人的個人資料",
			Data:    nil,
		})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("UpdateUserHandler bad request", "reason", "bind body", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "請求格式錯誤",
			Data:    nil,
		})
	}

	req.NickName = strings.TrimSpace(req.NickName)
	req.Description = strings.TrimSpace(req.Description)
	req.Avatar = strings.TrimSpace(req.Avatar)

	if req.NickName == "" {
		slog.Warn("UpdateUserHandler bad request", "reason", "empty nickName", "userId", userID)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "暱稱不可為空",
			Data:    nil,
		})
	}

	if err := validateAvatarURL(req.Avatar); err != nil {
		slog.Warn("UpdateUserHandler bad request", "reason", "invalid avatar", "userId", userID)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "大頭照 URL 格式不正確（需以 http:// 或 https:// 開頭）",
			Data:    nil,
		})
	}

	user, err := db.UpdateUser(db.Dbs, userID, req.NickName, req.Description, req.Avatar)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("UpdateUserHandler not found", "id", id)
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "找不到該使用者",
				Data:    nil,
			})
		}
		slog.Error("UpdateUser", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "更新失敗，請稍後再試",
			Data:    nil,
		})
	}

	auth, err := db.GetUserAuthByUserID(db.Dbs, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Warn("UpdateUserHandler user auth not found", "userId", userID)
			return c.Status(fiber.StatusUnauthorized).JSON(dto.TResponse[any]{
				Message: "找不到登入帳號資料",
				Data:    nil,
			})
		}
		slog.Error("GetUserAuthByUserID", "err", err, "userId", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "更新失敗，請稍後再試",
			Data:    nil,
		})
	}

	sessionUser := session.NewUser(user, auth.Username)
	if !session.SetUser(c, sessionUser) {
		slog.Error("UpdateUserHandler session not available", "userId", userID)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "更新失敗，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("UpdateUserHandler success", "userId", userID)
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[dto.UserProfileResponse]{
		Message: "ok",
		Data:    toUserProfileResponse(user),
	})
}

func validateAvatarURL(avatar string) error {
	if avatar == "" {
		return nil
	}
	if strings.HasPrefix(avatar, "http://") || strings.HasPrefix(avatar, "https://") {
		return nil
	}
	return errors.New("invalid avatar url")
}
