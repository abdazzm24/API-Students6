package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func listStudents(
	repo *StudentRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		query, err := parseListQuery(c)

		if err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"query tidak valid",
				map[string]string{
					"query": err.Error(),
				},
			)
		}

		ctx, cancel := reqCtx(c)
		defer cancel()

		result, total, err := repo.List(
			ctx,
			query.Search,
			query.IsActive,
			query.Sort,
			query.Order,
			query.Page,
			query.Limit,
		)

		if err != nil {
			return sendError(
				c,
				fiber.StatusInternalServerError,
				"gagal mengambil daftar student",
				nil,
			)
		}

		page, totalPages := calculatePagination(
			query.Page,
			query.Limit,
			total,
		)

		meta := &Meta{
			Page:       page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: totalPages,
		}

		return sendResponse(
			c,
			fiber.StatusOK,
			"daftar student berhasil diambil",
			result,
			meta,
			nil,
		)
	}
}

func getStudent(
	repo *StudentRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id, err := parseID(c)

		if err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"id tidak valid",
				map[string]string{
					"id": err.Error(),
				},
			)
		}

		ctx, cancel := reqCtx(c)
		defer cancel()

		student, err := repo.FindByID(ctx, id)

		if err != nil {
			return sendError(
				c,
				fiber.StatusInternalServerError,
				"gagal mengambil detail student",
				nil,
			)
		}

		if student == nil {
			return sendError(
				c,
				fiber.StatusNotFound,
				"student tidak ditemukan",
				nil,
			)
		}

		return sendSuccess(
			c,
			fiber.StatusOK,
			"student berhasil ditemukan",
			student,
		)
	}
}

func createStudent(
	repo *StudentRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		var req CreateStudentRequest

		if err := c.BodyParser(&req); err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"body bukan JSON yang sah",
				nil,
			)
		}

		validationErrors := validateCreateRequest(req)

		if len(validationErrors) > 0 {
			return sendError(
				c,
				fiber.StatusUnprocessableEntity,
				"validasi isi gagal",
				validationErrors,
			)
		}

		ctx, cancel := reqCtx(c)
		defer cancel()

		student, err := repo.Create(ctx, req)

		if err != nil {

			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "23505") ||
				strings.Contains(err.Error(), "unique") {

				return sendError(
					c,
					fiber.StatusConflict,
					"NIM sudah digunakan",
					map[string]string{
						"nim": "NIM harus unik",
					},
				)
			}

			return sendError(
				c,
				fiber.StatusInternalServerError,
				"gagal membuat student",
				nil,
			)
		}

		return sendCreated(
			c,
			"student berhasil ditambahkan",
			student,
			"/api/v1/students/"+strconv.Itoa(student.ID),
		)
	}
}

func replaceStudent(
	repo *StudentRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id, err := parseID(c)

		if err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"id tidak valid",
				map[string]string{
					"id": err.Error(),
				},
			)
		}

		var req ReplaceStudentRequest

		if err := c.BodyParser(&req); err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"body bukan JSON yang sah",
				nil,
			)
		}

		validationErrors := validateReplaceRequest(req)

		if len(validationErrors) > 0 {
			return sendError(
				c,
				fiber.StatusUnprocessableEntity,
				"validasi isi gagal",
				validationErrors,
			)
		}

		ctx, cancel := reqCtx(c)
		defer cancel()

		student, err := repo.Replace(ctx, id, req)

		if err != nil {

			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "23505") ||
				strings.Contains(err.Error(), "unique") {

				return sendError(
					c,
					fiber.StatusConflict,
					"NIM sudah digunakan",
					map[string]string{
						"nim": "NIM harus unik",
					},
				)
			}

			return sendError(
				c,
				fiber.StatusInternalServerError,
				"gagal memperbarui student",
				nil,
			)
		}

		if student == nil {
			return sendError(
				c,
				fiber.StatusNotFound,
				"student tidak ditemukan",
				nil,
			)
		}

		return sendSuccess(
			c,
			fiber.StatusOK,
			"student berhasil diganti",
			student,
		)
	}
}

func patchStudent(
	repo *StudentRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id, err := parseID(c)

		if err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"id tidak valid",
				map[string]string{
					"id": err.Error(),
				},
			)
		}

		var req PatchStudentRequest

		if err := c.BodyParser(&req); err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"body bukan JSON yang sah",
				nil,
			)
		}

		validationErrors := validatePatchRequest(req)

		if len(validationErrors) > 0 {
			return sendError(
				c,
				fiber.StatusUnprocessableEntity,
				"validasi isi gagal",
				validationErrors,
			)
		}

		ctx, cancel := reqCtx(c)
		defer cancel()

		student, err := repo.Patch(ctx, id, req)

		if err != nil {

			if strings.Contains(err.Error(), "duplicate key") ||
				strings.Contains(err.Error(), "23505") ||
				strings.Contains(err.Error(), "unique") {

				return sendError(
					c,
					fiber.StatusConflict,
					"NIM sudah digunakan",
					map[string]string{
						"nim": "NIM harus unik",
					},
				)
			}

			return sendError(
				c,
				fiber.StatusInternalServerError,
				"gagal memperbarui sebagian student",
				nil,
			)
		}

		if student == nil {
			return sendError(
				c,
				fiber.StatusNotFound,
				"student tidak ditemukan",
				nil,
			)
		}

		return sendSuccess(
			c,
			fiber.StatusOK,
			"student berhasil diperbarui sebagian",
			student,
		)
	}
}

func deleteStudent(
	repo *StudentRepository,
) fiber.Handler {

	return func(c *fiber.Ctx) error {

		id, err := parseID(c)

		if err != nil {
			return sendError(
				c,
				fiber.StatusBadRequest,
				"id tidak valid",
				map[string]string{
					"id": err.Error(),
				},
			)
		}

		ctx, cancel := reqCtx(c)
		defer cancel()

		deleted, err := repo.Delete(ctx, id)

		if err != nil {
			return sendError(
				c,
				fiber.StatusInternalServerError,
				"gagal menghapus student",
				nil,
			)
		}

		if !deleted {
			return sendError(
				c,
				fiber.StatusNotFound,
				"student tidak ditemukan",
				nil,
			)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
