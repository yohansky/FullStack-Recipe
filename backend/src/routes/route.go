package routes

import (
	"be-recipe/src/controllers"
	"be-recipe/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World, From Fiber Yohan")
	})
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	app.Get("/recipes", controllers.AllRecipes)
	app.Get("/recipe/:id", controllers.GetRecipe)

	app.Use(middleware.IsAuth)

	app.Get("/user", controllers.User)
	app.Post("/logout", controllers.Logout)

	app.Get("/users", controllers.AllUsers)
	app.Get("/user/:id", controllers.GetUser)
	app.Put("/user/:id", controllers.UpdateUser)
	app.Delete("/user/:id", controllers.DeleteUser)

	app.Post("/recipes", controllers.CreateRecipe)
	app.Put("/recipe/:id", controllers.UpdatePhotoRecipe)
	app.Delete("/recipe/:id", controllers.DeleteRecipe)
	app.Get("/recipes/user/:id", controllers.GetRecipesByUserId)
}
