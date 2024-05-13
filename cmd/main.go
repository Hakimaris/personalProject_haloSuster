package main

import (
	// "fmt"
	// "strings"
	// "time"
	"HaloSuster/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New()

	api := app.Group("/api") // /api
	v1 := api.Group("/v1")
	user:= v1.Group("/user")
	user.Get("/",handlers.GetUserHandler)
		nurse:= user.Group("/nurse")
		nurse.Get("/",handlers.GetNurseHandler)
		it:= user.Group("/it")
		it.Get("/",handlers.GetITHandler)

	log.Fatal(app.Listen(":8080"))
}