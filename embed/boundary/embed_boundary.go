package boundary

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/audryus/steganocc/embed/control"
	"github.com/audryus/steganocc/logger"
	"github.com/gofiber/fiber/v2"
)

func New(app *fiber.App, logger *logger.Log) {
	app.Get("/embed", getIndex)
	app.Post("/embed", postEncrypt)
}

func getIndex(c *fiber.Ctx) error {
	return index(c, fiber.Map{})
}

type ErrItem struct {
	Index int
	Msg   string
}

func index(c *fiber.Ctx, bind fiber.Map) error {
	bind["embed"] = true
	return c.Render("embed/boundary/embed", bind)
}

func postEncrypt(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["file"]
	phrase := c.FormValue("phrase")

	errs := make([]string, 0)

	if len(phrase) == 0 {
		errs = append(errs, "Missing enigmatic phrase.")
	}

	if len(files) == 0 {
		errs = append(errs, "Missing file.")
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

	c.Set("Content-Type", "image/png")
	filename := file.Filename
	if !strings.HasSuffix(filename, ".png") {
		filename = filename[:len(filename)-4]
		filename = filename + ".png"
	}
	c.Attachment(filename)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))

	buf, err := control.Embed(phrase, filename, f)
	if err != nil {
		return err
	}

	return c.SendStream(bytes.NewReader(buf.Bytes()), buf.Len())
}
