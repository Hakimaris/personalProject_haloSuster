package handlers

import "github.com/gofiber/fiber/v2"

func GetITHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "im it handler!",
	})
}
