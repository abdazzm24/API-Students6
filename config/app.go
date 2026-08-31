package config

import (
	"log/slog"
	
	"api-students/app/service"
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

	app := fiber.New()

	middleware.Register(
		app,
		logger,
	)

	route.Register(
		app,
		pool,
		studentService,
	)

	return app
}