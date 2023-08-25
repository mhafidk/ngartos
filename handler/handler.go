package handler

import "github.com/gofiber/fiber/v2"

func Check(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "ok",
		"message": "All is well",
		"data": nil,
	})
}