package middlewares

import (
	"e-commerce-fiber/utils"
	"encoding/json"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateRegisterUser(c *fiber.Ctx) error {
	validator := validator.New()

	type UserInput struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}

	var user UserInput

	bodyBytes := c.Body()

	if err := json.Unmarshal(bodyBytes, &user); err != nil {
		return utils.ErrorResponse(c, 400, "missing field")
	}

	if err := validator.Struct(&user); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}
	return c.Next()

}
