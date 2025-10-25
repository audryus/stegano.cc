package boundary

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

func Handle(app *fiber.App) {
	static := fiber.Static{
		Compress:      true,
		MaxAge:        86400,
		CacheDuration: time.Hour,
	}

	app.Static("/public/favicon", "./public/favicon", static)
	app.Static("/public/css", "./public/css", static)
	app.Static("/public/fonts", "./public/fonts", static)
	app.Static("/public/images", "./public/images", static)
	app.Static("/public/js", "./public/js", static)
	app.Static("/public/webfonts", "./public/webfonts", static)
}
