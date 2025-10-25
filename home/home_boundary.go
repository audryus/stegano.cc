package home

import (
	"github.com/audryus/steganocc/logger"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App, logger *logger.Log) {
	app.Get("/", index)
}

func index(c *fiber.Ctx) error {
	return c.Render("home/index", fiber.Map{})
}
