package handler

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

type updateExercise struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func CreateExercise(c *fiber.Ctx) error {
	db := database.DB.Db
	exercise := new(model.Exercise)

	err := c.BodyParser(exercise)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Something is wrong with the input data",
			"data":    err,
		})
	}

	slug := strings.ReplaceAll(exercise.Name, " ", "-")
	exercise.Slug = strings.ToLower(slug)

	err = db.Create(&exercise).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create exercise",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise created",
		"data":    nil,
	})
}

func GetSingleExercise(c *fiber.Ctx) error {
	db := database.DB.Db

	slug := c.Params("slug")

	var exercise model.Exercise

	db.Find(&exercise, "slug = ?", slug)
	if exercise.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "Exercise not found",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise found",
		"data": fiber.Map{
			"name":      exercise.Name,
			"content":   exercise.Content,
			"slug":      exercise.Slug,
			"createdAt": exercise.CreatedAt,
			"updatedAt": exercise.UpdatedAt,
		},
	})
}

func UpdateExercise(c *fiber.Ctx) error {
	db := database.DB.Db

	var exercise model.Exercise

	slug := c.Params("slug")

	db.Find(&exercise, "slug = ?", slug)
	if exercise.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "Exercise not found",
			"data":    nil,
		})
	}

	var updateexerciseData updateExercise
	err := c.BodyParser(&updateexerciseData)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Something is wrong with the input data",
			"data":    err,
		})
	}

	exercise.Name = updateexerciseData.Name
	slugUpdate := strings.ReplaceAll(exercise.Name, " ", "-")
	exercise.Slug = strings.ToLower(slugUpdate)
	err = db.Save(&exercise).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not update the exercise",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise updated",
		"data": fiber.Map{
			"name":      exercise.Name,
			"content":   exercise.Content,
			"createdAt": exercise.CreatedAt,
			"updatedAt": exercise.UpdatedAt,
		},
	})
}

func DeleteExercise(c *fiber.Ctx) error {
	db := database.DB.Db

	var exercise model.Exercise

	slug := c.Params("slug")

	db.Find(&exercise, "slug = ?", slug)
	if exercise.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "Exercise not found",
			"data":    nil,
		})
	}

	err := db.Delete(&exercise, "slug = ?", slug).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete exercise",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Exercise deleted",
		"data":    nil,
	})
}

func GetAllExercises(c *fiber.Ctx) error {
	db := database.DB.Db

	var exercises []model.Exercise

	db.Select("name", "id", "slug").Find(&exercises)
	if len(exercises) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Exercise not found",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "sucess",
		"message": "Topic Found",
		"data": fiber.Map{
			"exercises":  exercises,
			"total_data": len(exercises),
		},
	})
}

func GetAllTopicExercises(c *fiber.Ctx) error {
	db := database.DB.Db

	var exercises []model.Exercise
	var topic model.Topic

	topic_slug := c.Params("topic_slug")
	db.Find(&topic, "slug = ?", topic_slug)
	if topic.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "Topic not found",
			"data":    nil,
		})
	}

	db.Select("name", "id", "slug").Where("topic_id = ?", topic.ID).Find(&exercises)
	if len(exercises) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Exercise not found",
			"data":    nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "sucess",
		"message": "Topic Found",
		"data": fiber.Map{
			"exercises":  exercises,
			"total_data": len(exercises),
		},
	})
}
