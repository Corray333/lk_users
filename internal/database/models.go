package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	UID          primitive.ObjectID   `json:"uid" bson:"_id"`
	Name         string               `json:"name" bson:"name"`
	Username     string               `json:"username" bson:"username"`
	Email        string               `json:"email" bson:"email"`
	Posts        []primitive.ObjectID `json:"posts" bson:"posts"`
	Friends      []primitive.ObjectID `json:"friends" bson:"friends"`
	Password     string               `json:"password" bson:"password"`
	Groups       []primitive.ObjectID `json:"groups" bson:"groups"`
	Confirmation int                  `json:"confirmation" bson:"confirmation,omitempty"`
}
