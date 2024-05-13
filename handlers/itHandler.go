package handlers

import "github.com/gofiber/fiber/v2"

func GetITHandler(c *fiber.Ctx) error {
    return c.SendString("This is the IT Handler")
}
