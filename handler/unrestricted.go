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

func Check(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "ok",
		"message": "All is well",
		"data": nil,
	})
}

func Login(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Something is wrong with your input",
			"data": err,
		})
	}

	var credential string
	if user.Username == "" {
		credential = user.Email
	}
	if user.Email == "" {
		credential = user.Username
	}

	if credential == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Email or password could not be blank",
			"data": nil,
		})
	}

	userPassword := user.Password

	if user.Username == "" {
		db.Find(&user, "email = ?", user.Email)
	}
	if user.Email == "" {
		db.Find(&user, "username = ?", user.Username)
	}
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status": "error",
			"message": "User not found",
			"data": nil,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status": "error",
			"message": "Email and password don't match",
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
		"message": "User logged in",
		"data": fiber.Map{
			"token": t,
			"exp": exp,
		},
	})
}