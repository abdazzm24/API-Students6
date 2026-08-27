package main

import (
	"context"
	"log"
	"time"

	"api-students/config"
	"api-students/database"
	"api-students/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// 1. Load konfigurasi dari .env
	config.LoadEnv()

	// 2. Membuat koneksi database
	pool, err := database.NewPool(context.Background())

	if err != nil {
		log.Fatalf("database: %v", err)
	}

	defer pool.Close()

	// 3. Membuat repository
	studentRepository := repository.NewStudentRepository(pool)

	// 4. Membuat aplikasi Fiber
	app := fiber.New()

	// 5. Middleware
	app.Use(requestid.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// 6. Health check
	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(
			c.UserContext(),
			2*time.Second,
		)

		defer cancel()

		if err := pool.Ping(ctx); err != nil {
			return sendError(
				c,
				fiber.StatusServiceUnavailable,
				"database tidak dapat dihubungi",
				nil,
			)
		}

		return sendSuccess(
			c,
			fiber.StatusOK,
			"server dan database berjalan",
			map[string]string{
				"status": "ok",
			},
		)
	})

	// 7. API version
	api := app.Group("/api/v1")

	// 8. Student routes
	studentsAPI := api.Group("/students")

	studentsAPI.Get("/", listStudents(studentRepository))
	studentsAPI.Get("/:id", getStudent(studentRepository))
	studentsAPI.Post("/", requireJSON, createStudent(studentRepository))
	studentsAPI.Put("/:id", requireJSON, replaceStudent(studentRepository))
	studentsAPI.Patch("/:id", requireJSON, patchStudent(studentRepository))
	studentsAPI.Delete("/:id", deleteStudent(studentRepository))

	// 9. Endpoint tidak ditemukan
	app.Use(func(c *fiber.Ctx) error {
		return sendError(
			c,
			fiber.StatusNotFound,
			"endpoint tidak ditemukan",
			nil,
		)
	})

	// 10. Jalankan server
	log.Fatal(app.Listen(":3000"))
}
