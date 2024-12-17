package mongoTLC

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"mainService/internal/domain"
)

type IServiceRepository interface {
	AddService(userID string, service *domain.ApiService) (string, error)
	GetServiceByID(serviceID string) (*domain.ApiService, error)
	GetServicesByIDs(serviceIDs ...string) ([]*domain.ApiService, error)
	GetAllServices() ([]*domain.ApiService, error)
	DeleteService(userID, serviceID string) error
	SearchServices(queryString string) ([]*domain.ApiService, error)
}

type mongoServiceRepository struct {
	DB   *mongo.Database
	Coll *mongo.Collection
}

func NewMongoServiceRepository(db *mongo.Database) IServiceRepository {
	return &mongoServiceRepository{
		DB:   db,
		Coll: db.Collection("service"),
	}
}

func (repo *mongoServiceRepository) AddService(userID string, service *domain.ApiService) (string, error) {
	service.UserID = userID

	dbService, err := service.ToDB()
	if err != nil {
		return "", err
	}

	res, err := repo.Coll.InsertOne(context.TODO(), *dbService)
	if err != nil {
		return "", err
	}

	serviceDBRef := bson.M{
		"$ref": "service",
		"$id":  res.InsertedID,
	}

	upd := bson.M{
		"$push": bson.M{"services": serviceDBRef},
	}

	mongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return "", BAD_USER_ID
	}

	_, err = repo.DB.Collection("user").UpdateByID(context.TODO(), mongoID, upd)
	if err != nil {
		return "", err
	}

	serviceID, _ := res.InsertedID.(bson.ObjectID)

	return serviceID.Hex(), nil
}

func (repo *mongoServiceRepository) GetServiceByID(serviceID string) (*domain.ApiService, error) {
	mongoID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		return nil, BAD_SERVICE_ID
	}

	dbInfo := new(domain.DBService)
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

func (repo *mongoServiceRepository) GetServicesByIDs(serviceIDs ...string) ([]*domain.ApiService, error) {
	mongoIDs := make([]bson.ObjectID, len(serviceIDs))
	for i, id := range serviceIDs {
		mongoID, err := bson.ObjectIDFromHex(id)
		if err != nil {
			return nil, BAD_SERVICE_ID
		}

		mongoIDs[i] = mongoID
	}

	filter := bson.M{
		"_id": bson.M{"$in": mongoIDs},
	}

	cursor, err := repo.Coll.Find(context.TODO(), filter)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, NOT_FOUND
	} else if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var DBresults []*domain.DBService
	for cursor.Next(context.TODO()) {
		service := new(domain.DBService)

		if err = cursor.Decode(service); err != nil {
			return nil, err
		}

		DBresults = append(DBresults, service)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	var ApiResults []*domain.ApiService
	for _, DBserv := range DBresults {
		ApiServ, err := DBserv.ToApi()
		if err != nil {
			return nil, err
		}

		ApiResults = append(ApiResults, ApiServ)
	}

	return ApiResults, nil
}

func (repo *mongoServiceRepository) isUserServiceOwner(userID, serviceID string) (bool, error) {
	userMongoID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return false, BAD_USER_ID
	}

	serviceMongoID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		return false, BAD_PET_ID
	}

	filter := bson.M{
		"_id":       serviceMongoID,
		"owner.$id": userMongoID,
	}

	docCount, err := repo.Coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return docCount != 0, nil
}

func (repo *mongoServiceRepository) GetAllServices() ([]*domain.ApiService, error) {
	cursor, err := repo.Coll.Find(context.TODO(), bson.M{})
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, NOT_FOUND
	} else if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	var DBresults []*domain.DBService
	for cursor.Next(context.TODO()) {
		service := new(domain.DBService)

		if err = cursor.Decode(service); err != nil {
			return nil, err
		}

		DBresults = append(DBresults, service)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	var ApiResults []*domain.ApiService
	for _, DBserv := range DBresults {
		ApiServ, err := DBserv.ToApi()
		if err != nil {
			return nil, err
		}

		ApiResults = append(ApiResults, ApiServ)
	}

	return ApiResults, nil
}

func (repo *mongoServiceRepository) DeleteService(userID, serviceID string) error {
	isOwner, err := repo.isUserServiceOwner(userID, serviceID)
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

	serviceMongoID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		return BAD_PET_ID
	}

	delRes, err := repo.Coll.DeleteOne(context.TODO(), bson.M{"_id": serviceMongoID})
	if err != nil {
		return err
	}
	if delRes.DeletedCount == 0 {
		fmt.Println("PPPPP")
		return NOT_FOUND
	}

	petDBRef := bson.M{
		"$ref": "service",
		"$id":  serviceMongoID,
	}

	update := bson.M{
		"$pull": bson.M{"services": petDBRef},
	}

	updRes, err := repo.DB.Collection("user").UpdateByID(context.TODO(), userMongoID, update)
	if err != nil {
		return err
	}
	if updRes.MatchedCount == 0 {
		return NOT_FOUND
	}

	return nil
}

func (repo *mongoServiceRepository) SearchServices(queryString string) ([]*domain.ApiService, error) {
	filter := bson.M{
		"$text": bson.M{
			"$search": queryString,
		},
	}

	projection := bson.D{
		{"score", bson.M{"$meta": "textScore"}},
	}

	opt := options.Find().SetProjection(projection).SetSort(bson.D{{"score", bson.M{"$meta": "textScore"}}})
	cursor, err := repo.Coll.Find(context.TODO(), filter, opt)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, NOT_FOUND
	} else if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var rawResults []*domain.DBServiceSerachResult
	for cursor.Next(context.TODO()) {
		var curRes domain.DBServiceSerachResult

		if err = cursor.Decode(&curRes); err != nil {
			return nil, err
		}

		rawResults = append(rawResults, &curRes)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	var results []*domain.ApiService
	for _, res := range rawResults {
		apiServ, err := res.ToApiService()
		if err != nil {
			return nil, err
		}

		results = append(results, apiServ)
	}

	return results, nil
}
