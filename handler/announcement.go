package handler

import (
	"errors"
	"kurohelper-api/dto"
	"kurohelper-api/session"
	"log/slog"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"kurohelperservice/db"
)

func toAnnouncementResponse(item db.Announcement) dto.AnnouncementResponse {
	return dto.AnnouncementResponse{
		ID:        item.ID,
		Category:  item.Category,
		Title:     item.Title,
		Content:   item.Content,
		Icon:      item.Icon,
		Thumbnail: item.Thumbnail,
		Image:     item.Image,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func optionalTrimmedPtr(s *string) *string {
	if s == nil {
		return nil
	}
	v := strings.TrimSpace(*s)
	if v == "" {
		return nil
	}
	return &v
}

func sessionUserID(c fiber.Ctx) int {
	if me := session.LoadUser(c); me != nil {
		return me.ID
	}
	return 0
}

func GetAnnouncementsHandler(c fiber.Ctx) error {
	category := strings.TrimSpace(c.Query("category"))

	var (
		list []db.Announcement
		err  error
	)
	if category == "" {
		list, err = db.GetAllAnnouncements(db.Dbs)
	} else {
		list, err = db.GetAnnouncementsByCategory(db.Dbs, category)
	}
	if err != nil {
		slog.Error("GetAnnouncements", "err", err, "category", category)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	resp := make([]dto.AnnouncementResponse, 0, len(list))
	for _, item := range list {
		resp = append(resp, toAnnouncementResponse(item))
	}

	slog.Info("GetAnnouncementsHandler success", "category", category, "count", len(resp))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[[]dto.AnnouncementResponse]{
		Message: "ok",
		Data:    resp,
	})
}

func CreateAnnouncementHandler(c fiber.Ctx) error {
	var req dto.CreateAnnouncementRequest
	if err := c.Bind().Body(&req); err != nil {
		slog.Warn("CreateAnnouncementHandler bad request", "reason", "bind body", "err", err)
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "請求格式錯誤",
			Data:    nil,
		})
	}

	req.Category = strings.TrimSpace(req.Category)
	req.Title = strings.TrimSpace(req.Title)
	req.Content = strings.TrimSpace(req.Content)
	icon := optionalTrimmedPtr(req.Icon)
	thumbnail := optionalTrimmedPtr(req.Thumbnail)
	image := optionalTrimmedPtr(req.Image)

	if req.Category == "" || req.Title == "" || req.Content == "" {
		slog.Warn("CreateAnnouncementHandler bad request", "reason", "empty required fields")
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "category、title、content 不可為空",
			Data:    nil,
		})
	}

	item, err := db.CreateAnnouncement(db.Dbs, req.Category, req.Title, req.Content, icon, thumbnail, image)
	if err != nil {
		slog.Error("CreateAnnouncement", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("CreateAnnouncementHandler success", "id", item.ID, "userId", sessionUserID(c))
	return c.Status(fiber.StatusCreated).JSON(dto.TResponse[dto.AnnouncementResponse]{
		Message: "ok",
		Data:    toAnnouncementResponse(item),
	})
}

func DeleteAnnouncementHandler(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		slog.Warn("DeleteAnnouncementHandler bad request", "id", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(dto.TResponse[any]{
			Message: "無效的公告 ID",
			Data:    nil,
		})
	}

	if _, err := db.GetAnnouncementByID(db.Dbs, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(dto.TResponse[any]{
				Message: "找不到公告",
				Data:    nil,
			})
		}
		slog.Error("GetAnnouncementByID", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	if err := db.DeleteAnnouncement(db.Dbs, id); err != nil {
		slog.Error("DeleteAnnouncement", "err", err, "id", id)
		return c.Status(fiber.StatusInternalServerError).JSON(dto.TResponse[any]{
			Message: "發生錯誤，請稍後再試",
			Data:    nil,
		})
	}

	slog.Info("DeleteAnnouncementHandler success", "id", id, "userId", sessionUserID(c))
	return c.Status(fiber.StatusOK).JSON(dto.TResponse[any]{
		Message: "ok",
		Data:    nil,
	})
}
