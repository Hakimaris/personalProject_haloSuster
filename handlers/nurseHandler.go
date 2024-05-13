package handlers

import "github.com/gofiber/fiber/v2"

func GetNurseHandler(c *fiber.Ctx) error {
    return c.SendString("This is the Nurse Handler")
}
