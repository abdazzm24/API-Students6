package route

import (
	"api-students/app/service"
	"api-students/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Register(
	app *fiber.App,
	pool *pgxpool.Pool,
	studentService *service.StudentService,
) {

	api := app.Group("/api/v1")

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {

		if err := pool.Ping(c.UserContext()); err != nil {
			return c.Status(
				fiber.StatusServiceUnavailable,
			).JSON(fiber.Map{
				"success": false,
				"message": "database tidak dapat dihubungi",
			})
		}

		return c.Status(
			fiber.StatusOK,
		).JSON(fiber.Map{
			"success": true,
			"message": "server dan database berjalan",
			"data": fiber.Map{
				"status": "ok",
			},
		})
	})

	// Student routes
	students := api.Group("/students")

	students.Get(
		"/",
		studentService.List,
	)

	students.Get(
		"/:id",
		studentService.Get,
	)

	students.Post(
		"/",
		middleware.RequireJSON,
		studentService.Create,
	)

	students.Put(
		"/:id",
		middleware.RequireJSON,
		studentService.Replace,
	)

	students.Patch(
		"/:id",
		middleware.RequireJSON,
		studentService.Patch,
	)

	students.Delete(
		"/:id",
		studentService.Delete,
	)

	// Endpoint tidak ditemukan
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(
			fiber.StatusNotFound,
		).JSON(fiber.Map{
			"success": false,
			"message": "endpoint tidak ditemukan",
		})
	})
}