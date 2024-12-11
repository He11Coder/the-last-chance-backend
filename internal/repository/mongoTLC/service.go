package mongoTLC

import (

	//"go.mongodb.org/mongo-driver/bson/primitive"

	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"mainService/internal/domain"
)

type IServiceRepository interface {
	AddService(service *domain.ApiService) error
	GetServiceByID(serviceID string) (*domain.ApiService, error)
	GetServicesByIDs(serviceIDs ...string) ([]*domain.ApiService, error)
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

func (repo *mongoServiceRepository) AddService(service *domain.ApiService) error {
	dbService, err := service.ToDB()
	if err != nil {
		return err
	}

	res, err := repo.Coll.InsertOne(context.TODO(), *dbService)
	if err != nil {
		return err
	}

	fmt.Println(res.InsertedID)
	serviceDBRef := bson.M{
		"$ref": "service",
		"$id":  res.InsertedID,
	}

	upd := bson.M{
		"$push": bson.M{"services": serviceDBRef},
	}

	mongoID, err := bson.ObjectIDFromHex(service.UserID)
	if err != nil {
		return BAD_USER_ID
	}

	_, err = repo.DB.Collection("user").UpdateByID(context.TODO(), mongoID, upd)
	if err != nil {
		return err
	}

	return nil
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
