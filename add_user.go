package thelastchancebackend

import (
	"context"
	"fmt"
	"mainService/app"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserInfo struct {
	UserID    bson.ObjectID `json:"user_id,omitempty" bson:"_id,omitempty"`
	Username  string        `json:"username" bson:"name"`
	Contacts  string        `json:"contacts" bson:"contact"`
	UserImage string        `json:"user_image" bson:"avatar_url"`
}

func main() {
	db, err := app.GetMongo()
	if err != nil {
		fmt.Print(err)
	}
	defer db.Disconnect(context.TODO())

	c := db.Database("tlc").Collection("user")

	u := UserInfo{
		Username:  "Сергей Иванов",
		Contacts:  "+79831238497",
		UserImage: "/assets/avatars/sergeant.png",
	}

	res, err := c.InsertOne(context.TODO(), u)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println("Ins ID:", res.InsertedID)
}
