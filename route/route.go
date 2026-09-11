package route

import (
	"api-students/app/service"
	"api-students/helper"
	"api-students/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	Pool *pgxpool.Pool

	JWT *helper.JWTManager

	StudentService *service.StudentService

	AuthService *service.AuthService
}

func Register(
	app *fiber.App,
	deps Dependencies,
) {

	api := app.Group(
		"/api/v1",
	)

	// ========
	// PUBLIC
	// ========

	api.Get(
		"/health",
		func(c *fiber.Ctx) error {

			if err := deps.Pool.Ping(
				c.UserContext(),
			); err != nil {

				return helper.Fail(
					c,
					fiber.StatusServiceUnavailable,
					"database tidak dapat dihubungi",
					nil,
				)
			}

			return helper.Success(
				c,
				fiber.StatusOK,
				"server dan database berjalan",
				fiber.Map{
					"status": "ok",
				},
			)
		},
	)

	// ================
	// AUTHENTICATION
	// ================

	auth := api.Group(
		"/auth",
		middleware.RequireJSON,
	)

	// Register
	auth.Post(
		"/register",
		deps.AuthService.Register,
	)

	// Login + rate limiter
	auth.Post(
		"/login",
		middleware.LoginRateLimiter(),
		deps.AuthService.Login,
	)

	// Refresh token
	auth.Post(
		"/refresh",
		deps.AuthService.Refresh,
	)

	// Logout
	auth.Post(
		"/logout",
		deps.AuthService.Logout,
	)

	// Profile
	auth.Get(
		"/me",
		middleware.RequireAuth(
			deps.JWT,
		),
		deps.AuthService.Me,
	)

	// =======================
	// STUDENTS - WAJIB LOGIN
	// =======================

	students := api.Group(
		"/students",
		middleware.RequireAuth(
			deps.JWT,
		),
	)

	// GET /students
	students.Get(
		"/",
		deps.StudentService.List,
	)

	// GET /students/:id
	students.Get(
		"/:id",
		deps.StudentService.Get,
	)

	// POST /students
	students.Post(
		"/",
		middleware.RequireJSON,
		deps.StudentService.Create,
	)

	// PUT /students/:id
	students.Put(
		"/:id",
		middleware.RequireJSON,
		deps.StudentService.Replace,
	)

	// PATCH /students/:id
	students.Patch(
		"/:id",
		middleware.RequireJSON,
		deps.StudentService.Patch,
	)

	// DELETE /students/:id
	students.Delete(
		"/:id",
		deps.StudentService.Delete,
	)
}