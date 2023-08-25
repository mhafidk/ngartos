package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

type updateUser struct {
	Username string `json:"username"`
}

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

func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.Db

	id := c.Params("id")

	var user model.User

	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "User not found",
			"data": nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User found",
		"data": user,
	})
}

func UpdateUser(c *fiber.Ctx) error {
	db := database.DB.Db

	var user model.User

	id := c.Params("id")

	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "User not found",
			"data": nil,
		})
	}

	var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with the input data",
			"data": err,
		})
	}

	user.Username = updateUserData.Username
	db.Save(&user)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User is successfully updated",
		"data": user,
	})
}