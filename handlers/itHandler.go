package handlers

import (
	// "strings"
	// "time"
	"HaloSuster/db"
	"HaloSuster/helpers"
	"HaloSuster/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func UserLogin(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var loginResult models.UserModel

	if err := c.BodyParser(&loginResult); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	//Check if request is empty
	if loginResult.NIP == 0 || loginResult.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip or password is empty",
		})
	}

	// Check nip format
	if !helpers.ValidateNIP(loginResult.NIP) {
		// fmt.Println("nip exist")
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid",
		})
	}

	// Check if NIP exists
	var count int
	// err_nip := conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&count)
	err := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&count)
	// fmt.Println("nip exist success")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "nip not found",
		})
	}

	if strconv.FormatInt(loginResult.NIP, 10)[:3] != "615" {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid for it user",
		})
	}

	// get user data
	var dbpassword string
	err = conn.QueryRow("SELECT id, name, password FROM \"Users\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&loginResult.ID, &loginResult.Name, &dbpassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	// fmt.Println(dbpassword)

	// check password
	if !helpers.CheckPasswordHash(loginResult.Password, dbpassword) {
		return c.Status(400).JSON(fiber.Map{
			"message": "password is incorrect",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User logged in successfully",
		"data": struct {
			Id          string `json:"userId"`
			NIP         int64  `json:"nip"`
			Name        string `json:"name"`
			AccessToken string `json:"accessToken"`
		}{
			Id:          loginResult.ID,
			NIP:         loginResult.NIP,
			Name:        loginResult.Name,
			AccessToken: helpers.SignUserJWT(loginResult), // Add the appropriate value for AccessToken
		},
	})
}

func UserRegister(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var registerResult models.UserModel

	if err := c.BodyParser(&registerResult); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	//Check if request is empty
	if registerResult.NIP == 0 || registerResult.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip or password is empty",
		})
	}
	// Check nip format
	if !helpers.ValidateNIP(registerResult.NIP) {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid",
		})
	}
	// fmt.Println("nip format is valid")

	if strconv.FormatInt(registerResult.NIP, 10)[:3] != "615" {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid for it user",
		})
	}

	// Check if NIP already exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&count)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if count > 0 {
		return c.Status(409).JSON(fiber.Map{
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
	newPass, err := helpers.HashPassword(registerResult.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error hashing password",
		})
	}

	// insert data
	_, err = conn.Exec("INSERT INTO \"Users\" (nip, name, password) VALUES ($1, $2, $3)", registerResult.NIP, registerResult.Name, newPass)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// get inserted data
	err = conn.QueryRow("SELECT id FROM \"Users\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&registerResult.ID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": struct {
			Id          string `json:"userId"`
			NIP         int64  `json:"nip"`
			Name        string `json:"name"`
			AccessToken string `json:"accessToken"`
		}{
			Id:          registerResult.ID,
			NIP:         registerResult.NIP,
			Name:        registerResult.Name,
			AccessToken: helpers.SignUserJWT(registerResult), // Add the appropriate value for AccessToken
		},
	})
}
