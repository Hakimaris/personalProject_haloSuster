package handlers

import (
	"HaloSuster/db"
	"HaloSuster/helpers"
	"HaloSuster/models"
	"fmt"
	"log"
	"strconv"

	// "fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func MedicalAddPatient(c *fiber.Ctx) error {
	conn := db.CreateConn()

	var patientRequest models.PatientModel
	if err := c.BodyParser(&patientRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}
	fmt.Println("IdentityNumber:", patientRequest.IdentityNumber)
	fmt.Println("Name:", patientRequest.Name)
	fmt.Println("PhoneNumber:", patientRequest.PhoneNumber)
	fmt.Println("BirthDate:", patientRequest.BirthDate)
	fmt.Println("IdentityCardScanImg:", patientRequest.IdentityCardScanImg)
	fmt.Println("Gender:", patientRequest.Gender)

	if patientRequest.IdentityNumber == 0 || patientRequest.Name == "" || patientRequest.PhoneNumber == "" || patientRequest.BirthDate == "" || patientRequest.IdentityCardScanImg == "" || patientRequest.Gender == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "identity info not complete",
			// "data":		patientRequest.IdentityNumber,
			// "data2":		patientRequest.Name,
			// "data3":		patientRequest.PhoneNumber,
			// "data4":		patientRequest.BirthDate,
			// "data5":		patientRequest.IdentityCardScanImg,
			// "data6":		patientRequest.Gender,
		})
	}

	if !helpers.ValidateIdentity(patientRequest.IdentityNumber) {
		return c.Status(400).JSON(fiber.Map{
			"message": "identity number is invalid",
		})
	}

	if !helpers.ValidatePhoneNumber(patientRequest.PhoneNumber) {
		return c.Status(400).JSON(fiber.Map{
			"message": "phone number is invalid",
		})
	}

	if !helpers.ValidateURL(patientRequest.IdentityCardScanImg) {
		return c.Status(400).JSON(fiber.Map{
			"message": "identity card scan image is invalid",
		})
	}

	if !helpers.ValidateBirthDate(patientRequest.BirthDate) {
		return c.Status(400).JSON(fiber.Map{
			"message": "birth date is invalid",
		})
	}

	patientRequest.Gender = models.Gender(strings.ToLower(string(patientRequest.Gender)))
	if !helpers.ValidateGender(patientRequest.Gender) {
		return c.Status(400).JSON(fiber.Map{
			"message": "gender is invalid",
		})
	}

	// Check if identity number exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM patient WHERE \"identityNumber\" = $1 LIMIT 1", patientRequest.IdentityNumber).Scan(&count)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if count != 0 {
		return c.Status(409).JSON(fiber.Map{
			"message": "patient record already existed",
		})
	}

	// Insert the data
	_, err = conn.Exec("INSERT INTO \"patient\" (\"identityNumber\", name, \"phoneNumber\", \"birthDate\", gender, \"identityCardScanning\" ) VALUES ($1, $2, $3, $4, $5, $6)", patientRequest.IdentityNumber, patientRequest.Name, patientRequest.PhoneNumber, patientRequest.BirthDate, patientRequest.Gender, patientRequest.IdentityCardScanImg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "success",
	})
}

func MedicalAddRecord(c *fiber.Ctx) error {
	conn := db.CreateConn()
	userNip := c.Locals("userNip")
	userID := c.Locals("userId")

	var recordRequest models.RecordModel
	if err := c.BodyParser(&recordRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "error parsing body",
		})
	}

	if recordRequest.IdentityNumber == 0 || recordRequest.Symptoms == "" || recordRequest.Medications == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "record info is not complete",
		})
	}

	if !helpers.ValidateIdentity(recordRequest.IdentityNumber) {
		return c.Status(400).JSON(fiber.Map{
			"message": "identity is invalid",
		})
	}

	// Check if identity number exists
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM patient WHERE \"identityNumber\" = $1 LIMIT 1", recordRequest.IdentityNumber).Scan(&count)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if count == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "patient record not found",
		})
	}

	// Check if the user is authorized to add the record
	// fmt.Println(userNip)
	// fmt.Println(userID)
	err = conn.QueryRow("SELECT COUNT(*) FROM  \"Users\" WHERE nip = $1 AND id = $2 AND password IS NOT NULL LIMIT 1", userNip, userID).Scan(&count)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}
	if count == 0 {
		return c.Status(401).JSON(fiber.Map{
			"message": "user is not authorized to add record",
		})
	}

	// Insert the data
	fmt.Println("IdentityNumber:", recordRequest.IdentityNumber)
	fmt.Println("Symptoms:", recordRequest.Symptoms)
	fmt.Println("Medications:", recordRequest.Medications)
	fmt.Println("UserID:", userID)
	_, err = conn.Exec("INSERT INTO \"record\" (\"identityNumber\", symptoms, medications, \"createdBy\") VALUES ($1, $2, $3, $4)", recordRequest.IdentityNumber, recordRequest.Symptoms, recordRequest.Medications, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "record successfully added",
	})
}

func MedicalGetPatient(c *fiber.Ctx) error {
	conn := db.CreateConn() // Use the global db connection

	// Initialize query parameters
	identityNumber := c.Query("identityNumber")
	name := c.Query("name")
	phoneNumber := c.Query("phoneNumber")
	limit := c.Query("limit", "5")
	offset := c.Query("offset", "0")
	createdAt := c.Query("createdAt")

	// Build the base query
	query := "SELECT \"identityNumber\", \"phoneNumber\", name, \"birthDate\", gender, \"createdAt\" FROM patient WHERE 1=1"
	args := map[string]interface{}{}

	// Add filters based on the query parameters
	if identityNumber != "" {
		query += " AND id = :userId"
		args["userId"] = identityNumber
	}

	if name != "" {
		query += " AND LOWER(name) LIKE '%' || LOWER(:name) || '%'"
		args["name"] = name
	}

	if phoneNumber != "" {
		query += " AND CAST(nip AS text) LIKE '%' || :nip || '%'"
		args["nip"] = phoneNumber
	}

	if createdAt == "asc" || createdAt == "desc" {
		query += " ORDER BY \"createdAt\" " + createdAt
	} else {
		// Default order if not specified or invalid
		query += " ORDER BY \"createdAt\" DESC"
	}

	// Add limit and offset
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 5
	}
	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		offsetInt = 0
	}

	// Add limit and offset to the query string
	query += " LIMIT :limit OFFSET :offset"
	args["limit"] = limitInt
	args["offset"] = offsetInt

	namedQuery, err := conn.PrepareNamed(query)
	fmt.Print(namedQuery)
	if err != nil {
		log.Println("Failed to prepare the query:", err)
		return c.Status(500).SendString(err.Error())
	}

	// Define a slice to hold the results
	var patients []models.PatientModel

	// Execute the query and load the results into the patients slice
	err = namedQuery.Select(&patients, args)
	if err != nil {
		log.Println("Failed to execute the query:", err)
		return c.Status(500).SendString(err.Error())
	}

	// Return the results as JSON
	return c.Status(200).JSON(patients)
	// BELUM FORMATTING patients
}

func MedicalGetRecord(c *fiber.Ctx) error {
	// BELUM SEMUA
	return c.Status(200).JSON(fiber.Map{
		"message": "im it handler!",
	})

}
