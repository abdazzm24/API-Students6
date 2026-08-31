package helper

import (
	"api-students/app/model"

	"github.com/gofiber/fiber/v2"
)

func Success(
	c *fiber.Ctx,
	status int,
	message string,
	data any,
) error {

	return c.Status(status).JSON(
		model.WebResponse{
			Success: true,
			Message: message,
			Data:    data,
		},
	)
}

func Created(
	c *fiber.Ctx,
	message string,
	data any,
	location string,
) error {

	c.Set("Location", location)

	return c.Status(fiber.StatusCreated).JSON(
		model.WebResponse{
			Success: true,
			Message: message,
			Data:    data,
		},
	)
}

func Fail(
	c *fiber.Ctx,
	status int,
	message string,
	errors any,
) error {

	return c.Status(status).JSON(
		model.WebResponse{
			Success: false,
			Message: message,
			Errors:  errors,
		},
	)
}