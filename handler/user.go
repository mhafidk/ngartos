package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
)

type updateUser struct {
	Username string `json:"username"`
}

func GetCurrentUser(c *fiber.Ctx) error {
	db := database.DB.Db

	var user model.User

	currentUser := c.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db.Find(&user, "email = ?", email)
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
		"data": fiber.Map{
			"username": user.Username,
			"email": user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.Db

	id := c.Params("id")

	var user model.User

	currentUser := c.Locals("user").(*jwt.Token)
	claims := currentUser.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	db.Find(&user, "id = ?", id)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "User not found",
			"data": nil,
		})
	}

	if email != user.Email {
		return c.Status(404).JSON(fiber.Map{
			"status": "not found",
			"message": "User not found",
			"data": nil,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User found",
		"data": fiber.Map{
			"username": user.Username,
			"email": user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
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
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with the input data",
			"data": err,
		})
	}

	user.Username = updateUserData.Username
	db.Save(&user)

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User updated",
		"data": fiber.Map{
			"username": user.Username,
			"email": user.Email,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

func DeleteUser(c *fiber.Ctx) error {
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

	err := db.Delete(&user, "id = ?", id).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Failed to delete user",
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User deleted",
		"data": nil,
	})
}