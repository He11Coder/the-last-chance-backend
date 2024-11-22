package mongoTLC

import (
	"context"
	"errors"
	"fmt"

	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"mainService/internal/domain"
)

type IUserRepository interface {
	AddUser() error
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetAvatarPath(userID string) (string, error)
	//GetUserPets() error
}

type mongoUserRepository struct {
	DB   *mongo.Database
	Coll *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) IUserRepository {
	return &mongoUserRepository{
		DB:   db,
		Coll: db.Collection("user"),
	}
}

func (repo *mongoUserRepository) AddUser() error {
	return nil
}

func (repo *mongoUserRepository) GetUserInfo(userID string) (*domain.ApiUserInfo, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, BAD_USER_ID
	}

	dbInfo := new(domain.DBUserInfo)
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(dbInfo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, NOT_FOUND
	} else if err != nil {
		return nil, err
	}

	apiInfo, err := dbInfo.ToApi()
	if err != nil {
		return nil, err
	}

	return apiInfo, nil
}

func (repo *mongoUserRepository) AddPet(userID string, pet *domain.ApiPetInfo) error {
	pet_col := repo.DB.Collection("pet")

	res, err := pet_col.InsertOne(context.TODO(), *pet)
	if err != nil {
		return err
	}

	fmt.Println(res.InsertedID)
	petDBRef := bson.M{
		"$ref": "pet",
		"$id":  res.InsertedID,
	}

	upd := bson.M{
		"$push": bson.M{"pets": petDBRef},
	}

	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return BAD_USER_ID
	}

	_, err = repo.Coll.UpdateByID(context.TODO(), mongoID, upd)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mongoUserRepository) AddService(service *domain.DBService) error {
	//serv_col := repo.DB.Collection("service")

	//res, err := serv_col.InsertOne(context.TODO())
	return nil
}

func (repo *mongoUserRepository) GetAvatarPath(userID string) (string, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return "", BAD_USER_ID
	}

	var ava_url struct {
		avatar_url string `bson:"avatar_url"`
	}

	opt := options.FindOne().SetProjection(bson.M{"avatar_url": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&ava_url)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return ava_url.avatar_url, nil
}
