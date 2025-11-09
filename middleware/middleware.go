package middleware

import (
	"strconv"
	"time"

	"github.com/audryus/steganocc/config"
	"github.com/audryus/steganocc/logger"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/storage/bbolt/v2"
)

func Fiber(app *fiber.App, cfg config.Config, logger *logger.Log) {
	app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: logger.Engine(),
	}))
	app.Use(compress.New())
	app.Use(cache.New(cache.Config{
		Expiration:   24 * time.Hour,
		CacheControl: true,
	}))

	app.Use(requestid.New(requestid.Config{
		Header: "x-kong-request-id",
	}))

	app.Use(func(c *fiber.Ctx) error {
		err := c.Next()

		c.Response().Header.Set("X-Content-Type-Options", "nosniff")
		c.Response().Header.Set("X-Frame-Options", "DENY")
		c.Response().Header.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		c.Response().Header.Set("Content-Security-Policy", "default-src 'self'; img-src * 'self' data: https:; style-src 'self' 'unsafe-inline'; script-src 'self'")
		c.Response().Header.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Response().Header.Set("Permissions-Policy", "geolocation=(self), microphone=()")
		c.Response().Header.Set("X-XSS-Protection", "1; mode=block")

		return err
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     "GET, POST, OPTIONS",
		AllowOrigins:     cfg.Server.Addr,
	}))
	ratelimit, err := strconv.Atoi(cfg.App.RateLimit)
	if err != nil {
		logger.Error("Rate limit not valid.", err)
	}

	if cfg.App.Env != "local" {
		app.Use(limiter.New(limiter.Config{
			Max:        ratelimit,
			Expiration: 60 * time.Second,
			KeyGenerator: func(c *fiber.Ctx) string {
				if xf := c.Get("X-Forwarded-For"); xf != "" {
					return xf
				}
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				c.Set("Retry-After", "60")
				return c.SendStatus(fiber.StatusTooManyRequests)
			},
			SkipFailedRequests:     true,
			SkipSuccessfulRequests: false,
			LimiterMiddleware:      limiter.SlidingWindow{},

			Storage: bbolt.New(),
		}))
	}
}
