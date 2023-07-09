package users

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/Corray333/lk_users/internal/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Making temp user,
func NewUser(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano())

	user := database.User{}
	c.BodyParser(&user)
	user.UID = primitive.NewObjectID()
	user.Friends = []primitive.ObjectID{}
	user.Groups = []primitive.ObjectID{}
	user.Posts = []primitive.ObjectID{}
	user.Confirmation = rand.Intn(9000) + 1000

	if err := database.UsersDB.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Err(); err == nil {
		return c.SendStatus(http.StatusForbidden)
	}
	if err := database.UsersDB.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Err(); err == nil {
		return c.SendStatus(http.StatusForbidden)
	}

	auth := smtp.PlainAuth(
		"",
		"info@corray.ru",
		"QAZqaz555",
		"mail.netangels.ru",
	)
	err := smtp.SendMail(
		"mail.netangels.ru:25",
		auth,
		"info@corray.ru",
		[]string{user.Email},
		[]byte("From: info@corray.ru\nSubject: Lolkek\nYour confirmation code is: "+strconv.Itoa(user.Confirmation)),
	)
	if err != nil {
		fmt.Println(err)
	}

	_, err = database.UsersDB.InsertOne(context.TODO(), user)
	if err != nil {
		return c.SendStatus(http.StatusForbidden)
	}
	return c.SendStatus(http.StatusCreated)
}

func ConfirmUser(c *fiber.Ctx) error {
	type Request struct {
		Email        string `json:"email"`
		Confirmation int    `json:"confirmation"`
	}
	var req Request
	c.BodyParser(&req)
	var user database.User
	err := database.UsersDB.FindOne(context.TODO(), bson.D{{"email", req.Email}}).Decode(&user)
	fmt.Println(req.Email)
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}
	if user.Confirmation != req.Confirmation {
		return c.SendStatus(http.StatusLocked)
	}
	database.UsersDB.UpdateOne(context.TODO(), bson.D{{"email", req.Email}}, bson.D{{"$set", bson.D{{"confirmation", 0}}}})
	return c.SendStatus(http.StatusAccepted)
}
