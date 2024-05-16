package handlers

import (
	"HaloSuster/helpers"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UserGet(c *fiber.Ctx) error {
	metadata := c.Queries()
	// queries: userId, limit & offset, name, nip, role, createdAt

	// Default limit and offset values
	if limitStr := metadata["limit"]; limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Limit Param is invalid",
			})
		}
		offset := 0 // Assign the value to offset here
		// Use the limit and offset variables here
		_ = limit
		_ = offset
	}

	var err error
	offset := 0 // Declare the offset variable here
	if offsetStr := metadata["offset"]; offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"message": "Offset Param is invalid",
			})
		}
	}

	if role := metadata["role"]; role != "" {
		if !helpers.ValidateRoleRequest(role) {
			// Ignore invalid role
			delete(metadata, "role")
		}
	}

	if createdAt := metadata["createdAt"]; createdAt != "" {
		if !helpers.ValidateCreatedAtRequest(createdAt) {
			// Ignore invalid createdAt
			delete(metadata, "createdAt")
		}
	}

	// Example logic to fetch users based on queries
	users := []fiber.Map{}
	// Dummy response
	users = append(users, fiber.Map{
		"userId":    "123",
		"nip":       6152200102987,
		"name":      "John Doe",
		"createdAt": "2024-01-01T12:00:00Z",
	})

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"data":    users,
	})
}
