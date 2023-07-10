package users

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/Corray333/lk_users/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Secret []byte = []byte("IUBgioyt675T87g8oyfOb8g0Bg8vbVt98bg8Tt8B&92gigHhjgJHGjkgiohuIOBYguig")

// Creatind user
func NewUser(c *fiber.Ctx) error {
	rand.Seed(time.Now().UnixNano())

	// Parse request data
	user := database.User{}
	c.BodyParser(&user)
	user.UID = primitive.NewObjectID()
	user.Friends = []primitive.ObjectID{}
	user.Groups = []primitive.ObjectID{}
	user.Posts = []primitive.ObjectID{}
	user.Confirmation = rand.Intn(9000) + 1000
	user.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(user.Password)))

	// Check, if this user already exists
	if err := database.UsersDB.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Err(); err == nil {
		return c.SendStatus(http.StatusForbidden)
	}
	if err := database.UsersDB.FindOne(context.TODO(), bson.D{{"email", user.Email}}).Err(); err == nil {
		return c.SendStatus(http.StatusForbidden)
	}

	// Send email with confirmation code
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

	// Put user in database
	_, err = database.UsersDB.InsertOne(context.TODO(), user)
	if err != nil {
		return c.SendStatus(http.StatusForbidden)
	}
	return c.SendStatus(http.StatusCreated)
}

// Confirmation of registration with code from email
func ConfirmUser(c *fiber.Ctx) error {

	// Parse request data
	type Request struct {
		Email        string `json:"email"`
		Confirmation int    `json:"confirmation"`
	}
	var req Request
	c.BodyParser(&req)

	// Find user in database
	var user database.User
	err := database.UsersDB.FindOne(context.TODO(), bson.D{{"email", req.Email}}).Decode(&user)

	// Check if user exists
	if err != nil {
		return c.SendStatus(http.StatusNotFound)
	}

	// Check confirmation code
	if user.Confirmation != req.Confirmation {
		return c.SendStatus(http.StatusLocked)
	}

	// Update user confirmation
	database.UsersDB.UpdateOne(context.TODO(), bson.D{{"email", req.Email}}, bson.D{{"$set", bson.D{{"confirmation", 0}}}})
	return c.SendStatus(http.StatusAccepted)
}

func LogIn(c *fiber.Ctx) error {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Request
	c.BodyParser(&req)

	var user database.User
	err := database.UsersDB.FindOne(context.TODO(), bson.D{{"email", req.Email}}).Decode(&user)
	if err != nil {
		return c.SendStatus(http.StatusBadGateway)
	}
	if user.Confirmation != 0 {
		return c.SendString("Confirm your email!")
	}
	if user.Password == fmt.Sprintf("%x", sha256.Sum256([]byte(req.Password))) {
		// JWT signing
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
		claims["username"] = user.Username
		tokenStr, err := token.SignedString(Secret)
		if err != nil {
			return c.SendStatus(http.StatusBadGateway)
		}
		return c.SendString(tokenStr)
	}
	return c.SendStatus(http.StatusUnauthorized)
}

func Authorize(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	var req Request
	c.BodyParser(&req)

	// token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
	// 	_, ok := token.Method.(*jwt.SigningMethodHMAC)
	// 	if !ok {
	// 		c.SendStatus(http.StatusUnauthorized)
	// 	}
	// 	return "", nil
	// })
	// if err != nil {
	// 	return c.SendStatus(http.StatusUnauthorized)
	// }
	// if token.Valid {
	// 	return c.SendStatus(http.StatusAccepted)
	// }
	// return c.SendStatus(http.StatusUnauthorized)

	token, _ := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return Secret, nil
	})

	if claims, _ := token.Claims.(jwt.MapClaims); token.Valid {
		if claims["username"] == req.Username {
			return c.SendStatus(http.StatusAccepted)
		}
	}
	return c.SendStatus(http.StatusUnauthorized)
}
