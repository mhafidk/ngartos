package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

type updateTopic struct {
	Name string `json:"name"`
	Content string `json:"content"`
}

func CreateTopic(c *fiber.Ctx) error {
	db := database.DB.Db
	topic := new(model.Topic)

	err := c.BodyParser(topic)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with the input data",
			"data": err,
		})
	}

	err = db.Create(&topic).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Could not create a topic",
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Topic created",
		"data": nil,
	})
}

func GetSingleTopic(c *fiber.Ctx) error {
	db := database.DB.Db

	id := c.Params("id")

	var topic model.Topic

	db.Find(&topic, "id = ?", id)
	if topic.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "Topic not found",
			"data": nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User found",
		"data": fiber.Map{
			"name": topic.Name,
			"content": topic.Content,
			"createdAt": topic.CreatedAt,
			"updatedAt": topic.UpdatedAt,
		},
	})
}

func UpdateTopic(c *fiber.Ctx) error {
	db := database.DB.Db

	var topic model.Topic

	id := c.Params("id")

	db.Find(&topic, "id = ?", id)
	if topic.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "Topic not found",
			"data": nil,
		})
	}

	var updateTopicData updateTopic
	err := c.BodyParser(&updateTopicData)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with the input data",
			"data": err,
		})
	}

	topic.Name = updateTopicData.Name
	err = db.Save(&topic).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Could not update the topic",
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User updated",
		"data": fiber.Map{
			"name": topic.Name,
			"content": topic.Content,
			"createdAt": topic.CreatedAt,
			"updatedAt": topic.UpdatedAt,
		},
	})
}