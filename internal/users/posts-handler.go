package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Corray333/lk_users/internal/database"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SendPost(c *fiber.Ctx) error {

	// Parse request
	type Request struct {
		Users  []string `json:"users" bson:"users"`
		Post   string   `json:"post" bson:"post"`
		Author string   `json:"author" bson:"author"`
	}
	var req Request
	c.BodyParser(&req)

	// Give access to each user in the list
	for _, user := range req.Users {
		uid, _ := primitive.ObjectIDFromHex(user)
		fmt.Println(user)
		database.UsersDB.UpdateByID(context.TODO(), uid, bson.D{{"$push", bson.D{{"accessed", req.Post}}}})
	}

	// Add in "created" list of author
	database.UsersDB.UpdateOne(context.TODO(), bson.D{{"username", req.Author}}, bson.D{{"$push", bson.D{{"created", req.Post}}}})

	return c.SendStatus(http.StatusAccepted)
}
