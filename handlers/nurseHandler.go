package handlers

import (
	"HaloSuster/db"
	"HaloSuster/helpers"
	"HaloSuster/models"
	"database/sql"

	// "fmt"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func NurseLogin(c *fiber.Ctx) error {
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
	if strconv.FormatInt(loginResult.NIP, 10)[:3] != "303" {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid for it user",
		})
	}

	// Check if NIP exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&count)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if count == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip not found",
		})
	}

	// get user data
	var dbpassword sql.NullString
	err = conn.QueryRow("SELECT id, name, password FROM \"Users\" WHERE nip = $1 LIMIT 1", loginResult.NIP).Scan(&loginResult.ID, &loginResult.Name, &dbpassword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if !dbpassword.Valid {
		return c.Status(400).JSON(fiber.Map{
			"message": "User still has no access to the system",
		})
	}

	// check password
	if !helpers.CheckPasswordHash(loginResult.Password, dbpassword.String) {
		return c.Status(400).JSON(fiber.Map{
			"message": "password is incorrect",
		})
	}

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

	//Check whether the user exists
	// var count int
	// err_db := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 AND id = $2 LIMIT 1", userNip, userId).Scan(&count)
	// if err_db != nil {
	// 	return c.Status(500).JSON(fiber.Map{
	// 		"message": err_db,
	// 	})
	// }
	// if count == 0 {
	// 	return c.Status(401).JSON(fiber.Map{
	// 		"message": "User doesn't exist, please login properly",
	// 	})
	// }
}

func NurseRegister(c *fiber.Ctx) error {
	conn := db.CreateConn()
	// userNip := c.Locals("userNip")
	// userId := c.Locals("userId")

	var nurseModel models.UserModel
	if err := c.BodyParser(&nurseModel); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	if !helpers.ValidateNIP(nurseModel.NIP) {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid",
		})
	}

	if strconv.FormatInt(nurseModel.NIP, 10)[:3] != "303" {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid for nurse",
		})
	}

	if !helpers.ValidateURL(nurseModel.IdentityCardScanning) {
		return c.Status(400).JSON(fiber.Map{
			"message": "identityCardScanning format is invalid url",
		})
	}

	// Check if NIP exists
	var count int
	err_db := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 LIMIT 1", nurseModel.NIP).Scan(&count)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db,
		})
	}
	if count == 1 {
		return c.Status(409).JSON(fiber.Map{
			"message": "nip existed",
		})
	}

	// insert data
	_, err_db = conn.Exec("INSERT INTO \"Users\" (nip, name, \"identityCardScanning\") VALUES ($1, $2, $3)", nurseModel.NIP, nurseModel.Name, nurseModel.IdentityCardScanning)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db.Error(),
		})
	}

	// get inserted data
	err_data := conn.QueryRow("SELECT id FROM \"Users\" WHERE nip = $1 LIMIT 1", nurseModel.NIP).Scan(&nurseModel.ID)
	if err_data != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_data.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": struct {
			Id   string `json:"id"`
			NIP  int64  `json:"nip"`
			Name string `json:"name"`
			// AccessToken string `json:"access_token"`
		}{
			Id:   nurseModel.ID,
			NIP:  nurseModel.NIP,
			Name: nurseModel.Name,
			// AccessToken: helpers.SignUserJWT(nurseModel), // Add the appropriate value for AccessToken
		},
	})
}

func NursePut(c *fiber.Ctx) error {
	userId := c.Params("userId")
	conn := db.CreateConn()

	var nurseModel models.UserModel
	if err := c.BodyParser(&nurseModel); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	if !helpers.ValidateNIP(nurseModel.NIP) {
		return c.Status(400).JSON(fiber.Map{
			"message": "nip format is invalid",
		})
	}

	if strconv.FormatInt(nurseModel.NIP, 10)[:3] != "303" {
		return c.Status(404).JSON(fiber.Map{
			"message": "nip format is invalid for nurse",
		})
	}

	// Check if User exists
	var countUser int
	err_db := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE id = $1 LIMIT 1", userId).Scan(&countUser)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db,
		})
	}
	if countUser == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Check if NIP exists
	var count int
	err_db = conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE nip = $1 AND id != $2 LIMIT 1", nurseModel.NIP, userId).Scan(&count)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db,
		})
	}
	if count == 1 {
		return c.Status(409).JSON(fiber.Map{
			"message": "conflict nip, nip already existed",
		})
	}

	// update data
	_, err_db = conn.Exec("UPDATE \"Users\" SET nip = $1, name = $2 WHERE id = $3", nurseModel.NIP, nurseModel.Name, userId)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User updated successfully",
	})
}

func NurseDelete(c *fiber.Ctx) error {
	userId := c.Params("userId")
	conn := db.CreateConn()

	// Check if User exists
	var countUser int
	err := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE id = $1 LIMIT 1", userId).Scan(&countUser)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if countUser == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Check if User is nurse
	var nip string
	err = conn.QueryRow("SELECT nip FROM \"Users\" WHERE id = $1 LIMIT 1", userId).Scan(&nip)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	// fmt.Print(nip[:3])
	if nip[:3] != "303" {
		return c.Status(404).JSON(fiber.Map{
			"message": "User is not a nurse",
		})
	}

	// delete data
	_, err = conn.Exec("DELETE FROM \"Users\" WHERE id = $1", userId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func NurseAccess(c *fiber.Ctx) error {
	userId := c.Params("userId")
	conn := db.CreateConn()

	var nurseModel models.UserModel
	if err := c.BodyParser(&nurseModel); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	if !helpers.ValidatePassword(nurseModel.Password) {
		return c.Status(400).JSON(fiber.Map{
			"message": "password format is invalid",
		})
	}

	// Check if User exists
	var countUser int
	err_db := conn.QueryRow("SELECT COUNT(*) FROM \"Users\" WHERE id = $1 LIMIT 1", userId).Scan(&countUser)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db,
		})
	}
	if countUser == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Check if User is nurse
	var nip string
	err_db = conn.QueryRow("SELECT nip FROM \"Users\" WHERE id = $1 LIMIT 1", userId).Scan(&nip)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db,
		})
	}

	if nip[:3] != "303" {
		return c.Status(404).JSON(fiber.Map{
			"message": "User is not a nurse",
		})
	}

	//update data
	newPass, err := helpers.HashPassword(nurseModel.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "error hashing password",
		})
	}
	_, err_db = conn.Exec("UPDATE \"Users\" SET password = $1 WHERE id = $2", newPass, userId)
	if err_db != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err_db.Error(),
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "User data updated successfully",
	})
}
