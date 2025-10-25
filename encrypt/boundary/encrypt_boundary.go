package boundary

import (
	"github.com/audryus/steganocc/encrypt/control"
	"github.com/audryus/steganocc/logger"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App, logger *logger.Log) {
	app.Get("/encrypt", getIndex)
	app.Post("/encrypt", postEncrypt)
}

func getIndex(c *fiber.Ctx) error {
	return index(c, fiber.Map{})
}

type ErrItem struct {
	Index int
	Msg   string
}

func index(c *fiber.Ctx, bind fiber.Map) error {
	bind["encrypt"] = true
	return c.Render("encrypt/boundary/encrypt", bind)
}

func postEncrypt(c *fiber.Ctx) error {
	password := c.FormValue("password")
	phrase := c.FormValue("phrase")

	errs := make([]string, 0)

	if len(password) == 0 {
		errs = append(errs, "Missing password.")
	}
	if len(password) == 0 {
		errs = append(errs, "Missing something to encrypt.")
	}
	/* if len(files) == 0 {
		errs = append(errs, "Missing file.")
	} */

	items := make([]ErrItem, 0, len(errs))
	for i, e := range errs {
		items = append(items, ErrItem{Index: i + 1, Msg: e}) // +1 if you want 1-based
	}

	if len(items) > 0 {
		return index(c, fiber.Map{
			"Errors": items,
		})
	}

	encrypted, err := control.Encrypt(password, phrase)
	if err != nil {
		return index(c, fiber.Map{
			"Errors": "Something very wrong happened. Perhaps a bigger PNG would work.",
		})
	}

	return index(c, fiber.Map{
		"encrypted": encrypted,
	})
}
