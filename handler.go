package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var students = []Student{
	{
		ID:       1,
		NIM:      "11223344",
		Name:     "Abdullah Azzam",
		Grade:    85,
		IsActive: true,
	},
	{
		ID:       2,
		NIM:      "11223345",
		Name:     "Budi Santoso",
		Grade:    78,
		IsActive: true,
	},
	{
		ID:       3,
		NIM:      "11223346",
		Name:     "Citra Lestari",
		Grade:    92,
		IsActive: false,
	},
}

var nextStudentID = 4

func findStudentIndexByID(id int) int {
	for index, student := range students {
		if student.ID == id {
			return index
		}
	}

	return -1
}

func isNIMExists(nim string, exceptID int) bool {
	for _, student := range students {
		if student.NIM == nim && student.ID != exceptID {
			return true
		}
	}

	return false
}

func listStudents(c *fiber.Ctx) error {
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

	result := make([]Student, 0)

	for _, student := range students {
		if query.Search != "" {
			if !strings.Contains(
				strings.ToLower(student.Name),
				strings.ToLower(query.Search),
			) {
				continue
			}
		}

		if query.IsActive != nil {
			if student.IsActive != *query.IsActive {
				continue
			}
		}

		result = append(result, student)
	}

	sort.Slice(result, func(i, j int) bool {
		var less bool

		switch query.Sort {
		case "id":
			less = result[i].ID < result[j].ID

		case "nim":
			less = result[i].NIM < result[j].NIM

		case "name":
			less = strings.ToLower(result[i].Name) <
				strings.ToLower(result[j].Name)

		case "grade":
			less = result[i].Grade < result[j].Grade

		case "is_active":
			less = !result[i].IsActive &&
				result[j].IsActive
		}

		if query.Order == "desc" {
			return !less
		}

		return less
	})

	total := len(result)

	start := (query.Page - 1) * query.Limit

	if start >= total {
		result = []Student{}
	} else {
		end := start + query.Limit

		if end > total {
			end = total
		}

		result = result[start:end]
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

func getStudent(c *fiber.Ctx) error {
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

	index := findStudentIndexByID(id)

	if index == -1 {
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
		students[index],
	)
}

func createStudent(c *fiber.Ctx) error {
	if err := requireJSON(c); err != nil {
		return err
	}

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

	if isNIMExists(req.NIM, 0) {
		return sendError(
			c,
			fiber.StatusConflict,
			"NIM sudah digunakan",
			map[string]string{
				"nim": "NIM harus unik",
			},
		)
	}

	student := Student{
		ID:       nextStudentID,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}

	nextStudentID++

	students = append(students, student)

	return sendCreated(
		c,
		"student berhasil ditambahkan",
		student,
		"/api/v1/students/"+strconv.Itoa(student.ID),
	)
}

func replaceStudent(c *fiber.Ctx) error {
	if err := requireJSON(c); err != nil {
		return err
	}

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

	index := findStudentIndexByID(id)

	if index == -1 {
		return sendError(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
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

	if isNIMExists(req.NIM, id) {
		return sendError(
			c,
			fiber.StatusConflict,
			"NIM sudah digunakan",
			map[string]string{
				"nim": "NIM harus unik",
			},
		)
	}

	students[index] = Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	}

	return sendSuccess(
		c,
		fiber.StatusOK,
		"student berhasil diganti",
		students[index],
	)
}

func patchStudent(c *fiber.Ctx) error {
	if err := requireJSON(c); err != nil {
		return err
	}

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

	index := findStudentIndexByID(id)

	if index == -1 {
		return sendError(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
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

	if req.NIM != nil {
		if isNIMExists(*req.NIM, id) {
			return sendError(
				c,
				fiber.StatusConflict,
				"NIM sudah digunakan",
				map[string]string{
					"nim": "NIM harus unik",
				},
			)
		}

		students[index].NIM = *req.NIM
	}

	if req.Name != nil {
		students[index].Name = *req.Name
	}

	if req.Grade != nil {
		students[index].Grade = *req.Grade
	}

	if req.IsActive != nil {
		students[index].IsActive = *req.IsActive
	}

	return sendSuccess(
		c,
		fiber.StatusOK,
		"student berhasil diperbarui sebagian",
		students[index],
	)
}

func deleteStudent(c *fiber.Ctx) error {
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

	index := findStudentIndexByID(id)

	if index == -1 {
		return sendError(
			c,
			fiber.StatusNotFound,
			"student tidak ditemukan",
			nil,
		)
	}

	students = append(
		students[:index],
		students[index+1:]...,
	)

	return c.SendStatus(fiber.StatusNoContent)
}