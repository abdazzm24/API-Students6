package middleware

import (
	"github.com/gofiber/fiber/v2"

	"api-students/helper"
)

// RequirePermission memeriksa apakah user memiliki permission tertentu.
//
// Middleware ini harus dipasang setelah RequireAuth.
//
// Urutan:
//
// RequireAuth
//     ↓
// RequirePermission
func RequirePermission(
	permissions *helper.PermissionSet,
	requiredPermission string,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Mengambil identitas user dari context Fiber.
		currentUser, exists := helper.CurrentUser(c)

		// Jika identitas tidak ditemukan,
		// berarti user belum terautentikasi.
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"success": false,
					"message": "unauthorized",
				},
			)
		}

		// Memeriksa permission berdasarkan role.
		allowed := permissions.Can(
			currentUser.Role,
			requiredPermission,
		)

		if !allowed {
			return c.Status(fiber.StatusForbidden).JSON(
				fiber.Map{
					"success": false,
					"message": "forbidden",
				},
			)
		}

		return c.Next()
	}
}

// RequireRole membatasi endpoint berdasarkan role.
//
// Middleware ini bersifat opsional.
// Untuk pembatasan yang lebih fleksibel, gunakan
// RequirePermission.
func RequireRole(
	allowedRoles ...string,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		currentUser, exists := helper.CurrentUser(c)

		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{
					"success": false,
					"message": "unauthorized",
				},
			)
		}

		for _, role := range allowedRoles {
			if currentUser.Role == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(
			fiber.Map{
				"success": false,
				"message": "forbidden",
			},
		)
	}
}