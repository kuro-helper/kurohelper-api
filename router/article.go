package router

import (
	"kurohelper-api/handler"
	"kurohelper-api/middlware"

	"github.com/gofiber/fiber/v3"
)

// Article 文章資源
func ArticleRouter(apiGroup fiber.Router) {
	group := apiGroup.Group("/article")

	group.Get("/", middlware.TokenAuth(false), func(c fiber.Ctx) error {
		return handler.GetArticlesHandler(c)
	})

	group.Post("/", middlware.TokenAuth(true), middlware.SessionAuth(), func(c fiber.Ctx) error {
		return handler.CreateArticleHandler(c)
	})
}
