package service

import (
	"strconv"
	"strings"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"

	"github.com/gofiber/fiber/v2"
)

type StudentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{
		repo: repo,
	}
}

func (s *StudentService) List(c *fiber.Ctx) error {
	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	query, err := helper.ParseListQuery(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"query tidak valid",
			map[string]string{
				"query": err.Error(),
			},
		)
	}

	result, total, err := s.repo.List(
		ctx,
		query.Search,
		query.IsActive,
		query.Sort,
		query.Order,
		query.Page,
		query.Limit,
	)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengambil daftar student",
			nil,
		)
	}

	meta := &model.Meta{
		Page:       query.Page,
		Limit:      query.Limit,
		Total:      total,
		TotalPages: CalculateTotalPages(
			query.Limit,
			total,
		),
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"daftar student berhasil diambil",
		map[string]any{
			"students": result,
			"meta":     meta,
		},
	)
}


func (s *StudentService) Get(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			map[string]string{
				"id": err.Error(),
			},
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.repo.FindByID(ctx, id)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal mengambil detail student",
			nil,
		)
	}

	if student == nil {
		return helper.Fail(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"student berhasil ditemukan",
		student,
	)
}

func (s *StudentService) Create(c *fiber.Ctx) error {
	var req model.CreateStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body bukan JSON yang sah",
			nil,
		)
	}

	validationErrors := ValidateCreateRequest(req)

	if len(validationErrors) > 0 {
		return helper.Fail(
			c,
			fiber.StatusUnprocessableEntity,
			"validasi isi gagal",
			validationErrors,
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.repo.Create(ctx, req)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "23505") ||
			strings.Contains(err.Error(), "unique") {

			return helper.Fail(
				c,
				fiber.StatusConflict,
				"NIM sudah digunakan",
				map[string]string{
					"nim": "NIM harus unik",
				},
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal membuat student",
			nil,
		)
	}

	return helper.Created(
		c,
		"student berhasil ditambahkan",
		student,
		"/api/v1/students/"+strconv.Itoa(student.ID),
	)
}

func (s *StudentService) Replace(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			map[string]string{
				"id": err.Error(),
			},
		)
	}

	var req model.ReplaceStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body bukan JSON yang sah",
			nil,
		)
	}

	validationErrors := ValidateReplaceRequest(req)

	if len(validationErrors) > 0 {
		return helper.Fail(
			c,
			fiber.StatusUnprocessableEntity,
			"validasi isi gagal",
			validationErrors,
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.repo.Replace(ctx, id, req)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "23505") ||
			strings.Contains(err.Error(), "unique") {

			return helper.Fail(
				c,
				fiber.StatusConflict,
				"NIM sudah digunakan",
				map[string]string{
					"nim": "NIM harus unik",
				},
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal memperbarui student",
			nil,
		)
	}

	if student == nil {
		return helper.Fail(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"student berhasil diganti",
		student,
	)
}

func (s *StudentService) Patch(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			map[string]string{
				"id": err.Error(),
			},
		)
	}

	var req model.PatchStudentRequest

	if err := c.BodyParser(&req); err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"body bukan JSON yang sah",
			nil,
		)
	}

	validationErrors := ValidatePatchRequest(req)

	if len(validationErrors) > 0 {
		return helper.Fail(
			c,
			fiber.StatusUnprocessableEntity,
			"validasi isi gagal",
			validationErrors,
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	student, err := s.repo.Patch(ctx, id, req)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") ||
			strings.Contains(err.Error(), "23505") ||
			strings.Contains(err.Error(), "unique") {

			return helper.Fail(
				c,
				fiber.StatusConflict,
				"NIM sudah digunakan",
				map[string]string{
					"nim": "NIM harus unik",
				},
			)
		}

		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal memperbarui sebagian student",
			nil,
		)
	}

	if student == nil {
		return helper.Fail(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
		)
	}

	return helper.Success(
		c,
		fiber.StatusOK,
		"student berhasil diperbarui sebagian",
		student,
	)
}

func (s *StudentService) Delete(c *fiber.Ctx) error {
	id, err := helper.ParseID(c)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusBadRequest,
			"id tidak valid",
			map[string]string{
				"id": err.Error(),
			},
		)
	}

	ctx, cancel := helper.RequestContext(c)
	defer cancel()

	deleted, err := s.repo.Delete(ctx, id)

	if err != nil {
		return helper.Fail(
			c,
			fiber.StatusInternalServerError,
			"gagal menghapus student",
			nil,
		)
	}

	if !deleted {
		return helper.Fail(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}
