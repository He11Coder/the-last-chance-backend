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
	GetUserLoginByID(userID string) (string, error)
	CheckUser(cred *domain.LoginCredentials) (string, error)
	AddUser(newUser *domain.ApiUserInfo) (string, error)
	UpdateUser(userID string, updInfo *domain.ApiUserUpdate) error
	GetUserInfo(userID string) (*domain.ApiUserInfo, error)
	GetAvatarBytes(userID string) ([]byte, error)
	AddPet(userID string, pet *domain.ApiPetInfo) error
	DeletePet(userID, petID string) error
	UpdatePet(userID, petID string, updInfo *domain.ApiPetUpdate) error
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

func (repo *mongoUserRepository) GetUserLoginByID(userID string) (string, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return "", BAD_USER_ID
	}

	var login struct {
		Login string `bson:"login"`
	}

	opt := options.FindOne().SetProjection(bson.M{"login": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&login)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", NOT_FOUND
	} else if err != nil {
		return "", err
	}

	return login.Login, nil
}

func (repo *mongoUserRepository) CheckUser(cred *domain.LoginCredentials) (string, error) {
	var userCr domain.DBUserInfo

	opt := options.FindOne().SetProjection(bson.M{"hashed_password": 1, "salt": 1, "_id": 1})
	err := repo.Coll.FindOne(context.TODO(), bson.M{"login": cred.Username}, opt).Decode(&userCr)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return "", NOT_FOUND
	} else if err != nil {
		return "", err
	}

	isEqual := authUtils.ComparePasswordAndHash(cred.Password, userCr.Salt, userCr.HashedPassword)
	if !isEqual {
		return "", INCORRECT_CREDENTIALS
	}

	return userCr.UserID.Hex(), nil
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

func (repo *mongoUserRepository) UpdateUser(userID string, updInfo *domain.ApiUserUpdate) error {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return BAD_USER_ID
	}

	dbUpd, err := updInfo.ToDB()
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": dbUpd,
	}

	_, err = repo.Coll.UpdateByID(context.TODO(), mongoID, update)
	if err != nil {
		return err
	}

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

func (repo *mongoUserRepository) DeletePet(userID, petID string) error {
	isOwner, err := repo.isUserPetOwner(userID, petID)
	if err != nil {
		return err
	}

	if !isOwner {
		return ACCESS_DENIED
	}

	userMongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return BAD_USER_ID
	}

	petMongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return BAD_PET_ID
	}

	delRes, err := repo.DB.Collection("pet").DeleteOne(context.TODO(), bson.M{"_id": petMongoID})
	if err != nil {
		return err
	}
	if delRes.DeletedCount == 0 {
		return NOT_FOUND
	}

	petDBRef := bson.M{
		"$ref": "pet",
		"$id":  petMongoID,
	}

	update := bson.M{
		"$pull": bson.M{"pets": petDBRef},
	}

	updRes, err := repo.Coll.UpdateByID(context.TODO(), userMongoID, update)
	if err != nil {
		return err
	}
	if updRes.MatchedCount == 0 {
		return NOT_FOUND
	}

	return nil
}

func (repo *mongoUserRepository) UpdatePet(userID, petID string, updInfo *domain.ApiPetUpdate) error {
	isOwner, err := repo.isUserPetOwner(userID, petID)
	if err != nil {
		return err
	}

	if !isOwner {
		return ACCESS_DENIED
	}

	petMongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return BAD_PET_ID
	}

	dbUpd, err := updInfo.ToDB()
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": dbUpd,
	}

	_, err = repo.DB.Collection("pet").UpdateByID(context.TODO(), petMongoID, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mongoUserRepository) isUserPetOwner(userID, petID string) (bool, error) {
	userMongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return false, BAD_USER_ID
	}

	petMongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return false, BAD_PET_ID
	}

	filter := bson.M{
		"_id": userMongoID,
		"pets": bson.M{
			"$elemMatch": bson.M{
				"$id": petMongoID,
			},
		},
	}

	docCount, err := repo.Coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return docCount != 0, nil
}

func (repo *mongoUserRepository) GetAvatarBytes(userID string) ([]byte, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, BAD_USER_ID
	}

	var avaUrl struct {
		AvatarBytes []byte `bson:"avatar_url"`
	}

	opt := options.FindOne().SetProjection(bson.M{"avatar_url": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&avaUrl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []byte{}, nil
	} else if err != nil {
		return nil, err
	}

	return avaUrl.AvatarBytes, nil
}

func (repo *mongoUserRepository) GetUserPets(userID string) ([]string, error) {
	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, BAD_USER_ID
	}

	var userPets struct {
		PetIDs []bson.M `bson:"pets"`
	}

	opt := options.FindOne().SetProjection(bson.M{"pets": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&userPets)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	strPetIDs := make([]string, len(userPets.PetIDs))
	for i, pet := range userPets.PetIDs {
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
		ServiceIDs []bson.M `bson:"services"`
	}

	opt := options.FindOne().SetProjection(bson.M{"services": 1, "_id": 0})
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&userServices)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	strServiceIDs := make([]string, len(userServices.ServiceIDs))
	for i, service := range userServices.ServiceIDs {
		serviceID, ok := service["$id"].(bson.ObjectID)
		if !ok {
			return nil, serverErrors.CAST_ERROR
		}

		strServiceIDs[i] = serviceID.Hex()
	}

	return strServiceIDs, nil
}
