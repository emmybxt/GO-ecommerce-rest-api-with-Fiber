package routes

import (
	"e-commerce-fiber/controllers"
	"e-commerce-fiber/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("/api/v1")

	api.Get("/", controllers.HomePage)

	userApi := app.Group("/api/v1/users/auth")

	userApi.Post("/register", middlewares.ValidateRegisterUser, controllers.RegisterUser)
	userApi.Post("/login", controllers.Login)
	userApi.Get("/profile", middlewares.Authentication, controllers.UserProfile)

	//address routes
	addressApi := app.Group("/api/v1/address", middlewares.Authentication)
	addressApi.Put("/update", controllers.UpdateAddress)

	//product routes

	productApi := app.Group("/api/v1/product", middlewares.Authentication)

	productApi.Post("/", controllers.CreateProduct)

	productApi.Get("/list", controllers.GetAllProducts)

	productApi.Get("/:id", controllers.GetProduct)


	//carts Routes

	cartsApi := app.Group("/api/v1/carts", middlewares.Authentication)

	cartsApi.Post("/", controllers.AddProductToCart)

}
