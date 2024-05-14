package main

import (
	// "fmt"
	// "strings"
	// "time"
	"HaloSuster/handlers"
	"HaloSuster/db"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func init() {
	db.CreateConn()
}
func main() {
	app := fiber.New()

	api := app.Group("/v1")
	user := api.Group("/user")
	user.Get("/", handlers.GetUserHandler)
	// The request below doesnt require JWT
	nurse := user.Group("/nurse")
	nurse.Post("/login", handlers.NurseLogin)
	it := user.Group("/it")
	it.Get("/", handlers.GetITHandler)
	it.Post("/register", handlers.ItRegister)
	// it.Post("/login", handlers.ItLogin)
	// The request below requires JWT
	nurse.Get("/", handlers.GetNurseHandler)

	log.Fatal(app.Listen(":8080"))
}
