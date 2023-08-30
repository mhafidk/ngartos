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
	exercises := api.Group("/exercises")
	bookmarks := api.Group("/bookmarks")

	api.Get("/check", handler.Check)
	api.Post("/login", handler.Login)

	users.Get("/verify/:token", handler.VerifyEmail)
	users.Post("/", handler.CreateUser)
	users.Post("/forgot-password", handler.ForgotPassword)
	users.Put("/reset-password/:token", handler.ResetPassword)

	topics.Get("/:slug", handler.GetSingleTopic)
	topics.Get("/", handler.GetAllTopics)

	exercises.Get("/:slug", handler.GetSingleExercise)
	exercises.Get("/", handler.GetAllExercises)
	exercises.Get("/topic/:topic_slug", handler.GetAllTopicExercises)

	jwtSecretKey := config.Config("JWT_SECRET_KEY")
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(jwtSecretKey)},
	}))

	users.Get("/me", handler.GetCurrentUser)
	users.Get("/:id", handler.GetSingleUser)
	users.Put("/:id", handler.UpdateUser)
	users.Delete("/:id", handler.DeleteUser)

	topics.Post("/", handler.CreateTopic)
	topics.Put("/:slug", handler.UpdateTopic)
	topics.Delete("/:slug", handler.DeleteTopic)

	exercises.Post("/", handler.CreateExercise)
	exercises.Put("/:slug", handler.UpdateExercise)
	exercises.Delete("/:slug", handler.DeleteExercise)

	bookmarks.Post("/", handler.CreateBookmark)
	bookmarks.Get("/:exercise_id", handler.GetSingleBookmark)
	bookmarks.Get("/", handler.GetAllUserBookmarks)
	bookmarks.Delete("/:exercise_id", handler.DeleteBookmark)
}
