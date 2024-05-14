package handlers

import (
	// "strings"
	// "time"
	"HaloSuster/db"
	"HaloSuster/helpers"
	"HaloSuster/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func GetITHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "im it handler!",
	})
}

func ItRegister(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var registerResult models.ItModel

	if err := c.BodyParser(&registerResult); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	fmt.Println("parsing body success")
	// Check nip format
	if !helpers.ValidateUserNIP(registerResult.NIP) {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid",
		})
	}
	fmt.Println("nip format is valid")

	// Check if NIP already exists
	var count int
	err_nip := conn.QueryRow("SELECT COUNT(*) FROM \"admin\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&count)
	if err_nip != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error fetching data",
		})
	}
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip already exists",
		})
	}

	// check name format
	if !helpers.ValidateName(registerResult.Name) {
		return c.Status(400).JSON(fiber.Map{
			"message": "name format should be between 5-50 characters",
		})
	}

	// check password format
	if !helpers.ValidatePassword(registerResult.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "password format should be between 8-33 characters",
		})
	}

	// hash password
	newPass, err_psw := helpers.HashPassword(registerResult.Password)
	if err_psw != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error hashing password",
		})
	}

	// insert data
	_, err_db := conn.Exec("INSERT INTO admin (nip, name, password) VALUES ($1, $2, $3)", registerResult.NIP, registerResult.Name, newPass)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db.Error(),
		})
	}

	// get inserted data
	err_data := conn.QueryRow("SELECT id FROM \"admin\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&registerResult.Id)
	if err_data != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_data.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": struct {
			Id          string  `json:"id"`
			NIP         int64  `json:"nip"`
			Name        string `json:"name"`
			AccessToken string `json:"access_token"`
		}{
			Id:          registerResult.Id,
			NIP:         registerResult.NIP,
			Name:        registerResult.Name,
			AccessToken: helpers.SignAdminJWT(registerResult), // Add the appropriate value for AccessToken
		},
	})
}
