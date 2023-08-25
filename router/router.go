package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhafidk/ngartos/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	// users := api.Group("/users")

	api.Get("/check", handler.Check)
}