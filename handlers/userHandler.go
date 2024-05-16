package handlers

import (
	"HaloSuster/db"
	"HaloSuster/helpers"
	"HaloSuster/models"

	"github.com/gofiber/fiber/v2"
)

func UserGet(c *fiber.Ctx) error {
	userId := c.Query("userId", "")
	limit := c.Params("limit")
	name := c.Params("name", "")
	nip := c.Params("nip")
	role := c.Params("role", "")
	createdAt := c.Params("createdAt")

	conn := db.CreateConn()
	if !helpers.ValidateRoleRequest(role) {
		return c.Status(400).JSON(fiber.Map{
			"message": "role is invalid",
		})
	}
}