package config

import (
	"log/slog"

	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"
	"api-students/route"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewApp(
	pool *pgxpool.Pool,
	studentService *service.StudentService,
	logger *slog.Logger,
) *fiber.App {

	app := fiber.New(
		fiber.Config{
			AppName: GetEnv(
				"APP_NAME",
				"Praktikum Backend Lanjut",
			),

			ErrorHandler: newErrorHandler(
				logger,
			),

			// Membatasi ukuran request.
			BodyLimit: 1 * 1024 * 1024,
		},
	)

	middleware.Register(
		app,
		logger,
		GetEnv(
			"ALLOWED_ORIGINS",
			"",
		),
	)

	route.Register(
		app,
		pool,
		studentService,
	)

	app.Use(
		func(c *fiber.Ctx) error {

			return helper.Fail(
				c,
				fiber.StatusNotFound,
				"endpoint tidak ditemukan",
				nil,
			)
		},
	)

	return app
}

func newErrorHandler(
	logger *slog.Logger,
) fiber.ErrorHandler {

	return func(
		c *fiber.Ctx,
		err error,
	) error {

		status :=
			fiber.StatusInternalServerError

		message :=
			"terjadi error pada server"

		if e, ok :=
			err.(*fiber.Error); ok {

			status = e.Code
			message = e.Message
		}

		logger.Error(
			"unhandled_error",

			slog.String(
				"path",
				c.Path(),
			),

			slog.Int(
				"status",
				status,
			),

			slog.String(
				"error",
				err.Error(),
			),
		)

		return helper.Fail(
			c,
			status,
			message,
			nil,
		)
	}
}