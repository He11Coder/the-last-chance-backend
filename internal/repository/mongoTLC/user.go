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
	"mainService/pkg/authUtils"
	"mainService/pkg/serverErrors"
)

type IUserRepository interface {
	ValidateLogin(login string) error
	CheckUser(cred *domain.LoginCredentials) (string, error)
	AddUser(newUser *domain.ApiUserInfo) (string, error)
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetAvatarPath(userID string) (string, error)
	AddPet(userID string, pet *domain.ApiPetInfo) error
	GetUserPets(userID string) ([]string, error)
	GetUserServices(userID string) ([]string, error)
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

func (repo *mongoUserRepository) ValidateLogin(login string) error {
	if len(login) == 0 {
		return EMPTY_LOGIN
	}

	filter := bson.M{
		"login": login,
	}

	var result bson.M
	err := repo.Coll.FindOne(context.TODO(), filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil
	} else if err != nil {
		return err
	} else {
		return LOGIN_EXISTS
	}
}

func (repo *mongoUserRepository) CheckUser(cred *domain.LoginCredentials) (string, error) {
	var userCred struct {
		hashed_pass []byte        `bson:"hashed_password"`
		salt        []byte        `bson:"salt"`
		id          bson.ObjectID `bson:"_id"`
	}

	opt := options.FindOne().SetProjection(bson.M{"hashed_password": 1, "salt": 1, "_id": 1})
	err := repo.Coll.FindOne(context.TODO(), bson.M{"login": cred.Username}, opt).Decode(&userCred)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", NOT_FOUND
	} else if err != nil {
		return "", err
	}

	isEqual := authUtils.ComparePasswordAndHash(cred.Password, userCred.salt, userCred.hashed_pass)
	if !isEqual {
		return "", INCORRECT_CREDENTIALS
	}

	return userCred.id.Hex(), nil
}

func (repo *mongoUserRepository) AddUser(newUser *domain.ApiUserInfo) (string, error) {
	password := newUser.Password
	dbUser, err := newUser.ToDB()
	if err != nil {
		return "", err
	}

	hashedPass, salt, err := authUtils.GenerateHash(password)
	if err != nil {
		return "", err
	}

	dbUser.HashedPassword = hashedPass
	dbUser.Salt = salt

	res, err := repo.Coll.InsertOne(context.TODO(), *dbUser)
	if err != nil {
		return "", err
	}

	userID, _ := res.InsertedID.(bson.ObjectID)
	return userID.Hex(), nil
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

func (repo *mongoUserRepository) GetUserServices(userID string) ([]string, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, BAD_USER_ID
	}

	var userServices struct {
		serviceIDs []bson.M `bson:"services"`
	}

	opt := options.FindOne().SetProjection(bson.M{"services": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&userServices)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	strServiceIDs := make([]string, len(userServices.serviceIDs))
	for i, service := range userServices.serviceIDs {
		serviceID, ok := service["$id"].(bson.ObjectID)
		if !ok {
			return nil, serverErrors.CAST_ERROR
		}

		strServiceIDs[i] = serviceID.Hex()
	}

	return strServiceIDs, nil
}
