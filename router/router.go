package router

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/mhafidk/ngartos/config"
	"github.com/mhafidk/ngartos/handler"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	users := api.Group("/users")
	topics := api.Group("/topics")

	api.Get("/check", handler.Check)
	api.Post("/login", handler.Login)
	users.Put("/verify/:token", handler.VerifyEmail)
	users.Post("/", handler.CreateUser)

	jwtSecretKey := config.Config("JWT_SECRET_KEY")
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecretKey)},
	}))

	users.Get("/me", handler.GetCurrentUser)
	users.Get("/:id", handler.GetSingleUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)

	topics.Post("/", handler.CreateTopic)
	topics.Get("/:id", handler.GetSingleTopic)
	topics.Get("/", handler.GetAllTopic)
	topics.Put("/:id", handler.UpdateTopic)
	topics.Delete("/:id", handler.DeleteTopic)
}