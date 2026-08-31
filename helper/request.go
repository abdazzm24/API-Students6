package helper

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"api-students/app/model"

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

func RequestContext(
	c *fiber.Ctx,
) (context.Context, context.CancelFunc) {

	return context.WithTimeout(
		c.UserContext(),
		5*time.Second,
	)
}

func ParseID(c *fiber.Ctx) (int, error) {

	id, err := strconv.Atoi(
		c.Params("id"),
	)

	if err != nil {
		return 0, fmt.Errorf(
			"id harus berupa angka",
		)
	}

	return id, nil
}

func ParseListQuery(
	c *fiber.Ctx,
) (model.ListQuery, error) {

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

	search := strings.TrimSpace(
		c.Query("search"),
	)

	sort := strings.ToLower(
		strings.TrimSpace(
			c.Query("sort", "id"),
		),
	)

	if !allowedSortFields[sort] {
		return model.ListQuery{}, fmt.Errorf(
			"field sort '%s' tidak diperbolehkan",
			sort,
		)
	}

	order := strings.ToLower(
		strings.TrimSpace(
			c.Query("order", "asc"),
		),
	)

	if order != "asc" &&
		order != "desc" {

		return model.ListQuery{}, fmt.Errorf(
			"order harus asc atau desc",
		)
	}

	query := model.ListQuery{
		Page:   page,
		Limit:  limit,
		Search: search,
		Sort:   sort,
		Order:  order,
	}

	if active := c.Context().
		QueryArgs().
		Peek("is_active"); len(active) > 0 {

		value, err := strconv.ParseBool(
			string(active),
		)

		if err != nil {
			return model.ListQuery{}, fmt.Errorf(
				"is_active harus true atau false",
			)
		}

		query.IsActive = &value
	}

	return query, nil
}