package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const maxLimit = 100

var allowedSortFields = map[string]bool{
	"id":        true,
	"nim":       true,
	"name":      true,
	"grade":     true,
	"is_active": true,
}

func sendResponse(
	c *fiber.Ctx,
	status int,
	message string,
	data any,
	meta *Meta,
	errors any,
) error {
	return c.Status(status).JSON(WebResponse{
		Success: status >= 200 && status < 300,
		Message: message,
		Data:    data,
		Meta:    meta,
		Errors:  errors,
	})
}

func sendSuccess(
	c *fiber.Ctx,
	status int,
	message string,
	data any,
) error {
	return sendResponse(c, status, message, data, nil, nil)
}

func sendCreated(
	c *fiber.Ctx,
	message string,
	data any,
	location string,
) error {
	c.Set("Location", location)

	return sendResponse(
		c,
		fiber.StatusCreated,
		message,
		data,
		nil,
		nil,
	)
}

func sendError(
	c *fiber.Ctx,
	status int,
	message string,
	errors any,
) error {
	return sendResponse(
		c,
		status,
		message,
		nil,
		nil,
		errors,
	)
}

func requireJSON(c *fiber.Ctx) error {
	contentType := c.Get("Content-Type")

	if !strings.HasPrefix(
		strings.ToLower(contentType),
		"application/json",
	) {
		return sendError(
			c,
			fiber.StatusUnsupportedMediaType,
			"Content-Type harus application/json",
			nil,
		)
	}

	return nil
}

func parseID(c *fiber.Ctx) (int, error) {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return 0, fmt.Errorf("id harus berupa angka")
	}

	return id, nil
}

func parseListQuery(c *fiber.Ctx) (ListQuery, error) {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	if limit > maxLimit {
		limit = maxLimit
	}

	search := strings.TrimSpace(c.Query("search"))

	sort := strings.ToLower(
		strings.TrimSpace(c.Query("sort", "id")),
	)

	if !allowedSortFields[sort] {
		return ListQuery{}, fmt.Errorf(
			"field sort '%s' tidak diperbolehkan",
			sort,
		)
	}

	order := strings.ToLower(
		strings.TrimSpace(c.Query("order", "asc")),
	)

	if order != "asc" && order != "desc" {
		return ListQuery{}, fmt.Errorf(
			"order harus asc atau desc",
		)
	}

	query := ListQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
		Sort:   sort,
		Order:  order,
	}

	if active := c.Context().QueryArgs().Peek("is_active"); len(active) > 0 {
		value, err := strconv.ParseBool(string(active))

		if err != nil {
			return ListQuery{}, fmt.Errorf(
				"is_active harus true atau false",
			)
		}

		query.IsActive = &value
	}

	return query, nil
}

func validateCreateRequest(req CreateStudentRequest) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(req.NIM) == "" {
		errors["nim"] = "NIM wajib diisi"
	}

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "nama wajib diisi"
	}

	if req.Grade < 0 || req.Grade > 100 {
		errors["grade"] = "grade harus berada pada rentang 0 sampai 100"
	}

	return errors
}

func validateReplaceRequest(req ReplaceStudentRequest) map[string]string {
	errors := make(map[string]string)

	if strings.TrimSpace(req.NIM) == "" {
		errors["nim"] = "NIM wajib diisi"
	}

	if strings.TrimSpace(req.Name) == "" {
		errors["name"] = "nama wajib diisi"
	}

	if req.Grade < 0 || req.Grade > 100 {
		errors["grade"] = "grade harus berada pada rentang 0 sampai 100"
	}

	return errors
}

func validatePatchRequest(req PatchStudentRequest) map[string]string {
	errors := make(map[string]string)

	if req.NIM != nil && strings.TrimSpace(*req.NIM) == "" {
		errors["nim"] = "NIM tidak boleh kosong"
	}

	if req.Name != nil && strings.TrimSpace(*req.Name) == "" {
		errors["name"] = "nama tidak boleh kosong"
	}

	if req.Grade != nil && (*req.Grade < 0 || *req.Grade > 100) {
		errors["grade"] = "grade harus berada pada rentang 0 sampai 100"
	}

	return errors
}

func calculatePagination(
	page int,
	limit int,
	total int,
) (int, int) {
	totalPages := 0

	if total > 0 {
		totalPages = int(
			math.Ceil(float64(total) / float64(limit)),
		)
	}

	return page, totalPages
}
