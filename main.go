package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New()

	app.Use(requestid.New())

	app.Use(logger.New())

	app.Use(cors.New())

	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		return sendSuccess(
			c,
			fiber.StatusOK,
			"API is healthy",
			map[string]string{
				"status": "ok",
			},
		)
	})

	api := app.Group("/api/v1")

	studentsAPI := api.Group("/students")

	studentsAPI.Get("/", listStudents)
	studentsAPI.Get("/:id", getStudent)
	studentsAPI.Post("/", createStudent)
	studentsAPI.Put("/:id", replaceStudent)
	studentsAPI.Patch("/:id", patchStudent)
	studentsAPI.Delete("/:id", deleteStudent)

	app.Use(func(c *fiber.Ctx) error {
		return sendError(
			c,
			fiber.StatusNotFound,
			"endpoint tidak ditemukan",
			nil,
		)
	})

	app.Listen(":3000")
}
