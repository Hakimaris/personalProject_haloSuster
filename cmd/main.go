package main

import (
	// "fmt"
	// "strings"
	// "time"
	"HaloSuster/db"
	"HaloSuster/handlers"
	"HaloSuster/helpers"

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
	nurse := user.Group("/nurse", helpers.AuthMiddleware)
	nurse.Post("/login", handlers.NurseLogin)
	it := user.Group("/it")
	it.Get("/", handlers.GetUserHandler)
	it.Post("/register", handlers.UserRegister)
	it.Post("/login", handlers.UserLogin)
	// The request below requires JWT
	nurse.Get("/", handlers.GetNurseHandler)

	log.Fatal(app.Listen(":8080"))
}
