package handlers

import (
	// "fmt"
	// "strings"
	// "time"
	"HaloSuster/models"
	"HaloSuster/db"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/log"
)

func GetITHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "im it handler!",
	})
}

func ItRegister(c *fiber.Ctx) error {
	it := new(models.ItModel)
	if err := c.BodyParser(it); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}
	var registerResult models.ItModel
	conn :=db.CreateConn()

	return c.Status(200).JSON(fiber.Map{
		"message": "im it handler!",
	})
}
