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
	// medical := api.Group("/medical")
	user := api.Group("/user")
	it := user.Group("/it")
	nurse := user.Group("/nurse")

	// The request below doesnt require JWT
	it.Post("/register", handlers.UserRegister)
	it.Post("/login", handlers.UserLogin)
	nurse.Post("/login", handlers.NurseLogin)
	
	// The request below requires JWT for Admin
	// nurse := user.Group("/nurse", helpers.AuthITMiddleware)
	nurse.Post("/register", helpers.AuthITMiddleware, handlers.NurseRegister)
	nurse.Put("/:userId", helpers.AuthITMiddleware, handlers.NursePut)
	nurse.Delete("/:userId", helpers.AuthITMiddleware, handlers.NurseDelete)
	nurse.Post("/:userId/access", helpers.AuthITMiddleware, handlers.NurseAccess)
	user.Get("/", helpers.AuthITMiddleware, handlers.GetUser)
	// nurse.Put("/:id", helpers.AuthITMiddleware, handlers.NursePut)
	

	// The request below requires JWT for either nurse or IT role
	// medical.Post("/patient", helpers.AuthAllMiddleware, handlers.MedicalAdd)
	// medical.Get("/patient", helpers.AuthAllMiddleware, handlers.MedicalGet)
	// medical.Post("/record", helpers.AuthAllMiddleware, handlers.RecordAdd)
	// medical.Get("/record", helpers.AuthAllMiddleware, handlers.RecordGet)

	log.Fatal(app.Listen(":8080"))
}
