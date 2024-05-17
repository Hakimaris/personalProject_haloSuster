package handlers

import (
	"fmt"
	"log"
	"strconv"

	"HaloSuster/db"
	"HaloSuster/models"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jackc/pgx/v4/stdlib" // Import pgx driver for sqlx
	"github.com/jmoiron/sqlx"
)

// type User struct {
// 	ID        string `json:"id" db:"id"`
// 	Name      string `json:"name" db:"name"`
// 	NIP       string `json:"nip" db:"nip"`
// 	Role      string `json:"role" db:"role"`
// 	CreatedAt string `json:"createdAt" db:"created_at"`
// }

func GetUser(c *fiber.Ctx) error {
	conn := db.CreateConn() // Use the global db connection

	// Initialize query parameters
	userId := c.Query("userId")
	name := c.Query("name")
	nip := c.Query("nip")
	role := c.Query("role")
	limit := c.Query("limit", "5")
	offset := c.Query("offset", "0")
	createdAt := c.Query("createdAt")

	// Build the base query
	query := "SELECT id, name, nip, \"createdAt\" FROM \"Users\" WHERE 1=1"
	args := map[string]interface{}{}

	// Add filters based on the query parameters
	if userId != "" {
		query += " AND id = :userId"
		args["userId"] = userId
	}

	if name != "" {
		query += " AND LOWER(name) LIKE '%' || LOWER(:name) || '%'"
		args["name"] = name
	}

	if nip != "" {
		query += " AND CAST(nip AS text) LIKE '%' || :nip || '%'"
		args["nip"] = nip
	}

	if role != "" {
		switch role {
		case "it":
			query += " AND CAST(nip AS text) LIKE '615%'"
		case "nurse":
			query += " AND CAST(nip AS text) LIKE '303%'"
		default:
			// Invalid role, ignore the parameter
		}
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
	args["limit"] = limitInt
	args["offset"] = offsetInt

	// Add limit and offset to the query string
	query += " LIMIT :limit OFFSET :offset"

	// Execute the query
	var users models.UserModel
	namedQuery, namedArgs, err := sqlx.Named(query, args)
	if err != nil {
		log.Println("Failed to prepare the query:", err)
		return c.Status(500).SendString(err.Error())
	}
	fmt.Print(namedQuery)
	namedQuery = conn.Rebind(namedQuery)

	err = conn.Select(&users, namedQuery, namedArgs...)
	if err != nil {
		log.Println("Failed to execute the query:", err)
		return c.Status(500).SendString(err.Error())
	}

	// Return the results as JSON
	return c.JSON(users)
}