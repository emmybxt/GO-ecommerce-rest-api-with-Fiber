package utils

import "github.com/gofiber/fiber/v2"

func ErrorResponse(c *fiber.Ctx, statusCode int, errorMessage string) error {
	return c.Status(statusCode).JSON(fiber.Map{
		"success": false,
		"error":   errorMessage,
	})
}

func SuccessMessage(c *fiber.Ctx, successMessage string, data interface{}) error {
	return c.Status(200).JSON(fiber.Map{
		"success": false,
		"message": successMessage,
		"data":    data,
	})
}
