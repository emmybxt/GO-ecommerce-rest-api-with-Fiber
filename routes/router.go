package routes

import (
	"e-commerce-fiber/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/", controllers.HomePage)

	userApi := app.Group("/users/auth")

	userApi.Post("/register", controllers.RegisterUser)
	userApi.Post("/login", controllers.Login)

}
