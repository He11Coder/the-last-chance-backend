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
	"mainService/pkg/serverErrors"
)

type IUserRepository interface {
	AddUser() error
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetAvatarPath(userID string) (string, error)
	GetUserPets(userID string) ([]string, error)
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

	dbInfo, err := pet.ToDB()
	if err != nil {
		return err
	}

	res, err := pet_col.InsertOne(context.TODO(), *dbInfo)
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

	var avaUrl struct {
		avatarUrl string `bson:"avatar_url"`
	}

	opt := options.FindOne().SetProjection(bson.M{"avatar_url": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&avaUrl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return avaUrl.avatarUrl, nil
}

func (repo *mongoUserRepository) GetUserPets(userID string) ([]string, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, BAD_USER_ID
	}

	var userPets struct {
		petIDs []bson.M `bson:"pets"`
	}

	opt := options.FindOne().SetProjection(bson.M{"pets": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&userPets)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	strPetIDs := make([]string, len(userPets.petIDs))
	for i, pet := range userPets.petIDs {
		petID, ok := pet["$id"].(bson.ObjectID)
		if !ok {
			return nil, serverErrors.CAST_ERROR
		}

		strPetIDs[i] = petID.Hex()
	}

	return strPetIDs, nil
}
