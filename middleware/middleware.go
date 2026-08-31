package middleware

import (
	"log/slog"
	"time"

	"api-students/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func RequireJSON(c *fiber.Ctx) error {
	if !c.Is("json") {
		return helper.Fail(
			c,
			fiber.StatusUnsupportedMediaType,
			"Content-Type harus application/json",
			nil,
		)
	}

	return c.Next()
}

func Register(
	app *fiber.App,
	logger *slog.Logger,
) {

	app.Use(requestid.New())

	app.Use(cors.New())

	app.Use(RequestLogger(logger))
}

func RequestLogger(
	logger *slog.Logger,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		start := time.Now()

		err := c.Next()

		logger.Info(
			"http_request",

			slog.String(
				"request_id",
				c.Get("X-Request-ID"),
			),

			slog.String(
				"method",
				c.Method(),
			),

			slog.String(
				"path",
				c.Path(),
			),

			slog.Int(
				"status",
				c.Response().StatusCode(),
			),

			slog.Duration(
				"duration",
				time.Since(start),
			),
		)

		return err
	}
}