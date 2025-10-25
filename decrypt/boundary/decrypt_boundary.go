package boundary

import (
	"github.com/audryus/steganocc/decrypt/control"
	"github.com/audryus/steganocc/logger"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App, logger *logger.Log) {
	app.Get("/decrypt", getIndex)
	app.Post("/decrypt", postDecrypt)
}

func getIndex(c *fiber.Ctx) error {
	return index(c, fiber.Map{})
}

func index(c *fiber.Ctx, bind fiber.Map) error {
	bind["decrypt"] = true
	return c.Render("decrypt/boundary/decrypt", bind)
}

type ErrItem struct {
	Index int
	Msg   string
}

func postDecrypt(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["file"]
	message := c.FormValue("message")
	password := c.FormValue("password")

	errs := make([]string, 0)

	if len(password) == 0 {
		errs = append(errs, "Missing password.")
	}

	if len(files) == 0 && len(message) == 0 {
		errs = append(errs, "Nothing to decrypt. Add a file or a encrypted message.")
	}

	items := make([]ErrItem, 0, len(errs))
	for i, e := range errs {
		items = append(items, ErrItem{Index: i + 1, Msg: e}) // +1 if you want 1-based
	}

	if len(items) > 0 {
		return index(c, fiber.Map{
			"Errors": items,
		})
	}

	file := files[0]

	f, err := file.Open()
	if err != nil {
		return err
	}

	decryptedMessage, decryptedFile, err := control.Decrypt(password, message, f)
	if err != nil {
		return index(c, fiber.Map{
			"Errors": "Error trying to decode the file.",
		})
	}

	return index(c, fiber.Map{
		"decryptedMessage": decryptedMessage,
		"decryptedFile":    decryptedFile,
	})
}
