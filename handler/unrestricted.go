package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mhafidk/ngartos/config"
	"github.com/mhafidk/ngartos/database"
	"github.com/mhafidk/ngartos/model"
	"github.com/mhafidk/ngartos/utils"
	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
)

func Check(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status":  "ok",
		"message": "All is well",
		"data":    nil,
	})
}

func Login(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Something is wrong with your input",
			"data":    err,
		})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email or password could not be blank",
			"data":    nil,
		})
	}

	userPassword := user.Password

	db.Find(&user, "email = ? OR username = ?", user.Email, user.Email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    nil,
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userPassword)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email and password don't match",
			"data":    err,
		})
	}

	if !user.Verified {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email hasn't been verified, please verify your email first!",
			"data":    nil,
		})
	}

	exp := time.Now().Add(time.Hour * 72).Unix()
	claims := jwt.MapClaims{
		"email": user.Email,
		"admin": false,
		"exp":   exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecretKey := config.Config("JWT_SECRET_KEY")
	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "There is something wrong",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User logged in",
		"data": fiber.Map{
			"token": t,
			"exp":   exp,
		},
	})
}

func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Something is wrong with the input data",
			"data":    err,
		})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Email or password could not be blank",
			"data":    nil,
		})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "There is something wrong",
			"data":    err,
		})
	}

	user.Password = string(hash)

	verificationToken := randstr.Hex(16)
	user.VerificationToken = verificationToken
	user.Verified = false

	err = db.Create(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not create user",
			"data":    err,
		})
	}

	utils.SendVerificationEmail(verificationToken, user.Email)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User created",
		"data":    nil,
	})
}

func VerifyEmail(c *fiber.Ctx) error {
	db := database.DB.Db

	token := c.Params("token")

	var user model.User

	db.Find(&user, "verification_token = ?", token)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "not found",
			"message": "User not found",
			"data":    nil,
		})
	}

	verifyAt := time.Now()
	user.VerifyAt = &verifyAt
	user.Verified = true
	user.VerificationToken = ""

	err := db.Save(&user).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Could not verify the user",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User verified",
		"data":    nil,
	})
}
