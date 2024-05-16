package handlers

import "github.com/gofiber/fiber/v2"

func GetNurseHandler(c *fiber.Ctx) error {
		userNip := c.Locals("userNip")
		userId := c.Locals("userId")
    return c.JSON(fiber.Map{
			"message": "me nursehandler",
			"userNip": userNip,
			"userId": userId,
		})
}

func NurseLogin(c *fiber.Ctx) error {
    return c.SendString("This is the Nurse Login Handler")
}


