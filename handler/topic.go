package handler

import (
	"strings"

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

	slug := strings.ReplaceAll(topic.Name, " ", "-")
	topic.Slug = strings.ToLower(slug)

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
	var topicChildren []model.Topic

	db.Find(&topic, "id = ?", id)
	if topic.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "Topic not found",
			"data": nil,
		})
	}

	db.Select("name", "id", "parent_id").Find(&topicChildren, "parent_id = ?", topic.ID)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Topic found",
		"data": fiber.Map{
			"name": topic.Name,
			"content": topic.Content,
			"createdAt": topic.CreatedAt,
			"updatedAt": topic.UpdatedAt,
			"children": topicChildren,
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
	slug := strings.ReplaceAll(topic.Name, " ", "-")
	topic.Slug = strings.ToLower(slug)
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
		"message": "Topic updated",
		"data": fiber.Map{
			"name": topic.Name,
			"content": topic.Content,
			"createdAt": topic.CreatedAt,
			"updatedAt": topic.UpdatedAt,
		},
	})
}

func DeleteTopic(c *fiber.Ctx) error {
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

	err := db.Delete(&topic, "id = ?", id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to delete topic",
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "Topic deleted",
		"data": nil,
	})
}

func GetAllTopics(c *fiber.Ctx) error {
	db := database.DB.Db

	var topics []model.Topic

	db.Select("name", "id", "parent_id").Find(&topics)
	if len(topics) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "Topic not found",
			"data": nil,
		})
	}

	var mainTopics []model.Topic
	groupedTopics := make(map[string][]model.Topic)
	for _, topic := range topics {
		if topic.ParentID == nil {
			mainTopics = append(mainTopics, topic)
		} else {
			groupedTopics[topic.ParentID.String()] = append(groupedTopics[topic.ParentID.String()], topic)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "sucess",
		"message": "Topic Found",
		"data": fiber.Map{
			"main_topics": mainTopics,
			"grouped_topics": groupedTopics,
			"total_topics": len(topics),
		},
	})
}