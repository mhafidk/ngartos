package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/config"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
	"golang.org/x/crypto/bcrypt"
)

type updateUser struct {
	Username string `json:"username"`
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with the input data",
			"data": err,
		})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Email or password could not be blank",
			"data": nil,
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "There is something wrong",
			"data": err,
		})
	}

	user.Password = string(hash)

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "Could not create a user",
			"data": err,
		})
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	claims := jwt.MapClaims{
		"email": user.Email,
		"admin": false,
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecretKey := config.Config("JWT_SECRET_KEY")
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "error",
			"message": "There is something wrong",
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"message": "User created",
		"data": fiber.Map{
			"token": t,
			"exp": exp,
		},
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
		"data": user,
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