package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/audryus/steganocc/config"
	decrypt "github.com/audryus/steganocc/decrypt/boundary"
	embed "github.com/audryus/steganocc/embed/boundary"
	encrypt "github.com/audryus/steganocc/encrypt/boundary"
	"github.com/audryus/steganocc/home"
	"github.com/audryus/steganocc/logger"
	"github.com/audryus/steganocc/middleware"
	public "github.com/audryus/steganocc/public"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"go.uber.org/fx"
)

func Register() fx.Option {
	return fx.Options(
		fx.Provide(
			logger.New,
			config.New,
		),
	)
}

func Routes() fx.Option {
	return fx.Options(
		fx.Invoke(middleware.Fiber),
		fx.Invoke(public.Handle),
		fx.Invoke(home.New),
		fx.Invoke(encrypt.New),
		fx.Invoke(decrypt.New),
		fx.Invoke(embed.New),
	)
}

func main() {
	app := fx.New(
		fx.Provide(CreateServer),
		Register(),
		Routes(),
		fx.Invoke(func(lifecycle fx.Lifecycle, logger *logger.Log, app *fiber.App, cfg config.Config) {
			lifecycle.Append(fx.Hook{
				OnStart: func(context.Context) error {
					logger.Info("Initializing server ...")
					go app.Listen(cfg.Http.Addr)
					return nil
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Shuting down ...")
					return app.Shutdown()
				},
			})
		}),
	)

	app.Run()
}

func CreateServer(cfg config.Config) *fiber.App {
	engine := django.NewFileSystem(http.Dir("./"), ".html")

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		ServerHeader: cfg.Server.Header,
		AppName:      cfg.App.Name,
		Views:        engine,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		BodyLimit:    10 * 1024 * 1024, // 10MB in bytes
	})

	return app
}
