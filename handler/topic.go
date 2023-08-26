package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

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