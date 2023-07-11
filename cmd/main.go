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

	app.Post("/new-user", users.NewUser)
	app.Post("/confirm-user", users.ConfirmUser)
	app.Post("/log-in", users.LogIn)
	app.Post("/authorize", users.Authorize)
	app.Post("/send-post", users.SendPost)

	log.Fatal(app.Listen("127.0.0.1:3000"))
}
