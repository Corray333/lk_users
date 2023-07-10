package main

import (
	"log"

	"github.com/Corray333/lk_users/internal/database"
	"github.com/Corray333/lk_users/internal/users"
	"github.com/gofiber/fiber/v2"
)

func main() {
	database.Connect()
	app := fiber.New()

	app.Put("/new-user", users.NewUser)
	app.Put("/confirm-user", users.ConfirmUser)
	app.Put("/log-in", users.LogIn)
	app.Put("/authorize", users.Authorize)

	log.Fatal(app.Listen("127.0.0.1:3000"))
}
