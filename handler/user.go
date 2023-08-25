package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with the input data",
			"data": err,
		})
	}

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Could not create a user",
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User is successfully created",
		"data": user,
	})
}