package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// Announcement 公告資源
func AnnouncementRouter(apiGroup fiber.Router) {
	group := apiGroup.Group("/announcement")

	group.Get("/", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetAnnouncementsHandler(c)
	})

	// 最高驗證：特權 Token + 登入 Session + DB 角色為開發者/站主
	group.Post("/", middlware.TokenAuth(true), middlware.SessionAuth(), middlware.RequireStaffRole(), func(c fiber.Ctx) error {
		return handler.CreateAnnouncementHandler(c)
	})

	group.Delete("/:id", middlware.TokenAuth(true), middlware.SessionAuth(), middlware.RequireStaffRole(), func(c fiber.Ctx) error {
		return handler.DeleteAnnouncementHandler(c)
	})
}
