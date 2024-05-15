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

func GetUserHandler(c *fiber.Ctx) error {
	return c.Status(200).JSON(fiber.Map{
		"message": "im it handler!",
	})
}

func UserLogin(c *fiber.Ctx) error {
	conn := db.CreateConn()
	var loginResult models.UserModel

	if err := c.BodyParser(&loginResult); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}
	fmt.Println("parsing body success")

	// Check nip format
	if !helpers.ValidateUserNIP(loginResult.NIP) {
		fmt.Println("nip exist")
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid",
		})
	}
	fmt.Println("nip success")

	// Check if NIP exists
	var count int
	// err_nip := conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&count)
	err_nip := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&count)
	fmt.Println("nip exist success")
	if err_nip != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_nip,
		})
	}
	if count == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip not found",
		})
	}

	// get user data
	var dbpassword string
	err_data := conn.QueryRow("SELECT id, name, password FROM \"Users\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&loginResult.ID, &loginResult.Name, &dbpassword)
	if err_data != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_data.Error(),
		})
	}
	fmt.Println(dbpassword)

	// check password
	if !helpers.CheckPasswordHash(loginResult.Password, dbpassword) {
		return c.Status(400).JSON(fiber.Map{
			"message": "password is incorrect",
		})
	}

	//   err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	// if err != nil {
	//     c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
	//     return
	// }
	// fmt.Println("parsing body success")

	return c.Status(200).JSON(fiber.Map{
		"message": "User logged in successfully",
		"data": struct {
			Id          string `json:"id"`
			NIP         int64  `json:"nip"`
			Name        string `json:"name"`
			AccessToken string `json:"access_token"`
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
	// err_nip := conn.QueryRow("SELECT COUNT(*) FROM \"user\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&count)
	err_nip := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&count)
	// err_nip := conn.QueryRow("select * from user")
	if err_nip != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_nip,
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
	_, err_db := conn.Exec("INSERT INTO \"Users\" (nip, name, password) VALUES ($1, $2, $3)", registerResult.NIP, registerResult.Name, newPass)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db.Error(),
		})
	}

	// get inserted data
	err_data := conn.QueryRow("SELECT id FROM \"Users\" WHERE nip = $1 LIMIT 1", registerResult.NIP).Scan(&registerResult.ID)
	if err_data != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_data.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": struct {
			Id          string `json:"id"`
			NIP         int64  `json:"nip"`
			Name        string `json:"name"`
			AccessToken string `json:"access_token"`
		}{
			Id:          registerResult.ID,
			NIP:         registerResult.NIP,
			Name:        registerResult.Name,
			AccessToken: helpers.SignUserJWT(registerResult), // Add the appropriate value for AccessToken
		},
	})
}
