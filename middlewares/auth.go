package middlewares

import (
	"e-commerce-fiber/utils"

	"github.com/gofiber/fiber/v2"
)

func Authentication(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	token := c.Cookies("jwt")

	if authHeader == "" || token == "" {
		return utils.ErrorResponse(c, 400, "Authorization Headers cannot be empty")
	}

	id,email,userType, err := utils.VerifyToken(token)

	if err!= nil {
		return utils.ErrorResponse(c, 400, "invalid token")
	}

	// set user id and email to context
	c.Locals("id", id)
	c.Locals("email", email)
	c.Locals("userType", userType)
	return c.Next()

}
