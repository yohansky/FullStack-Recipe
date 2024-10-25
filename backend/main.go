package main

import (
	"be-recipe/routes"
	"be-recipe/src/config"

	"github.com/gofiber/fiber/v2"
	"github.com/subosito/gotenv"
)

func main() {
	gotenv.Load()
	config.InitDB()

	app := fiber.New()
	routes.Router(app)
	app.Listen(":8080")
}
