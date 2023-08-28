package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

func CreateBookmark(c *fiber.Ctx) error {
	db := database.DB.Db
	bookmark := new(model.Bookmark)

	var user model.User

	currentUser := c.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db.Find(&user, "email = ?", email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "User not found",
			"data":    nil,
		})
	}

	bookmark.UserID = user.ID

	err := c.BodyParser(bookmark)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Something is wrong with the input data",
			"data":    err,
		})
	}

	err = db.Create(&bookmark).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create bookmark",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark created",
		"data":    nil,
	})
}

func GetSingleBookmark(c *fiber.Ctx) error {
	db := database.DB.Db

	exerciseID := c.Params("exercise_id")

	var bookmark model.Bookmark
	var user model.User

	currentUser := c.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db.Find(&user, "email = ?", email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "User not found",
			"data":    nil,
		})
	}

	var bookmarked bool
	var message string

	db.Find(&bookmark, "exercise_id = ? AND user_id = ?", exerciseID, user.ID)
	if bookmark.ID == uuid.Nil {
		bookmarked = false
		message = "Bookmarked not found"
	} else {
		bookmarked = true
		message = "Bookmarked found"
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": message,
		"data": fiber.Map{
			"bookmarked": bookmarked,
		},
	})
}

func DeleteBookmark(c *fiber.Ctx) error {
	db := database.DB.Db

	var bookmark model.Bookmark
	var user model.User

	exerciseID := c.Params("exercise_id")

	currentUser := c.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db.Find(&user, "email = ?", email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "User not found",
			"data":    nil,
		})
	}

	db.Find(&bookmark, "exercise_id = ? AND user_id = ?", exerciseID, user.ID)
	if bookmark.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "Bookmark not found",
			"data":    nil,
		})
	}

	err := db.Unscoped().Delete(&bookmark).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete exercise",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookmark deleted",
		"data":    nil,
	})
}

func GetAllUserBookmarks(c *fiber.Ctx) error {
	db := database.DB.Db

	var exercises []model.Exercise
	var user model.User

	currentUser := c.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db.Preload("Bookmarks").Find(&user, "email = ?", email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "User not found",
			"data":    nil,
		})
	}

	var exerciseIDs []string
	for _, bookmark := range user.Bookmarks {
		exerciseIDs = append(exerciseIDs, bookmark.ExerciseID.String())
	}

	db.Select("id", "name", "slug").Find(&exercises, exerciseIDs)

	return c.Status(200).JSON(fiber.Map{
		"status":  "sucess",
		"message": "Bookmark Found",
		"data": fiber.Map{
			"bookmarks":  exercises,
			"total_data": len(exercises),
		},
	})
}
