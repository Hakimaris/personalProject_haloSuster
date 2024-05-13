package handlers

import "github.com/gofiber/fiber/v2"

func GetUserHandler(c *fiber.Ctx) error {
    return c.SendString("This is the User Handler")
}
