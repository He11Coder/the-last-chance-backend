package mongoTLC

import (
	"context"
	"errors"

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
	SearchServices(queryString string, filters *domain.ServiceFilter) ([]*domain.ApiService, error)
	GetServiceIDsWithAnimals(animalList []string) ([]string, error)
}

type mongoServiceRepository struct {
	DB          *mongo.Database
	ServiceColl *mongo.Collection
	AnimalColl  *mongo.Collection
}

func NewMongoServiceRepository(db *mongo.Database) IServiceRepository {
	return &mongoServiceRepository{
		DB:          db,
		ServiceColl: db.Collection("service"),
		AnimalColl:  db.Collection("animal"),
	}
}

func (repo *mongoServiceRepository) AddService(userID string, service *domain.ApiService) (string, error) {
	service.UserID = userID

	dbService, err := service.ToDB()
	if err != nil {
		return "", err
	}

	res, err := repo.ServiceColl.InsertOne(context.TODO(), *dbService)
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
	err = repo.ServiceColl.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(dbInfo)
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

	cursor, err := repo.ServiceColl.Find(context.TODO(), filter)
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

	docCount, err := repo.ServiceColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	return docCount != 0, nil
}

func (repo *mongoServiceRepository) GetAllServices() ([]*domain.ApiService, error) {
	cursor, err := repo.ServiceColl.Find(context.TODO(), bson.M{})
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

	delRes, err := repo.ServiceColl.DeleteOne(context.TODO(), bson.M{"_id": serviceMongoID})
	if err != nil {
		return err
	}
	if delRes.DeletedCount == 0 {
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

func (repo *mongoServiceRepository) GetServiceIDsWithAnimals(animalList []string) ([]string, error) {
	filter := bson.M{
		"type_of_animal": bson.M{
			"$in": animalList,
		},
	}

	opt := options.Find().SetProjection(bson.M{"services.$id": 1, "_id": 0})
	cursor, err := repo.AnimalColl.Find(context.TODO(), filter, opt)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	serviceIDs := []string{}
	for _, res := range results {
		IDs := res["services"].(bson.A)
		for _, rawId := range IDs {
			castedID := rawId.(bson.D)
			for _, elem := range castedID {
				if elem.Key == "$id" {
					objID := elem.Value.(bson.ObjectID)
					serviceIDs = append(serviceIDs, objID.Hex())
				}
			}
		}
	}

	return serviceIDs, nil
}

func (repo *mongoServiceRepository) SearchServices(queryString string, filters *domain.ServiceFilter) ([]*domain.ApiService, error) {
	filter := bson.M{}

	if filters.MaxPrice == 0 && filters.MinPrice > 0 {
		filter["price"] = bson.M{"$gte": filters.MinPrice}
	} else if filters.MinPrice == 0 && filters.MaxPrice > 0 {
		filter["price"] = bson.M{"$lte": filters.MaxPrice}
	} else if filters.MinPrice > 0 && filters.MaxPrice > filters.MinPrice {
		filter["price"] = bson.M{"$gte": filters.MinPrice, "$lte": filters.MaxPrice}
	} else if filters.MinPrice > 0 && filters.MaxPrice == filters.MinPrice {
		filter["price"] = bson.M{"$eq": filters.MinPrice}
	}

	projection := bson.D{}
	sortRule := bson.D{}
	if queryString != "" {
		filter["$text"] = bson.M{"$search": queryString}

		projection = bson.D{
			{"score", bson.M{"$meta": "textScore"}},
		}
		sortRule = bson.D{{"score", bson.M{"$meta": "textScore"}}}
	}

	opt := options.Find().SetProjection(projection).SetSort(sortRule)
	cursor, err := repo.ServiceColl.Find(context.TODO(), filter, opt)
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

	results := []*domain.ApiService{}
	for _, res := range rawResults {
		apiServ, err := res.ToApiService()
		if err != nil {
			return nil, err
		}

		results = append(results, apiServ)
	}

	animalFilteredRes := []*domain.ApiService{}
	if len(filters.Animals) != 0 {
		IDsWithAnimals, err := repo.GetServiceIDsWithAnimals(filters.Animals)
		if err != nil {
			return nil, err
		}

		IDsToLeave := make(map[string]struct{}, len(IDsWithAnimals))
		for _, id := range IDsWithAnimals {
			IDsToLeave[id] = struct{}{}
		}

		for _, serv := range results {
			if _, exists := IDsToLeave[serv.ServiceID]; exists {
				animalFilteredRes = append(animalFilteredRes, serv)
			}
		}

		return animalFilteredRes, nil
	}

	return results, nil
}
