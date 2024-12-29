package mongoTLC

import (
	"context"
	"errors"
	"strings"

	//"fmt"
	"mainService/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IPetRepository interface {
	GetPetInfo(petID string) (*domain.ApiPetInfo, error)
	GetAvatarBytes(petID string) ([]byte, error)
	IncrementAnimal(typeOfAnimal string, serviceID string) error
	DecrementAnimal(typeOfAnimal string, serviceID string) error
	GetTopAnimals(top int64) ([]string, error)
}

type mongoPetRepository struct {
	DB         *mongo.Database
	PetColl    *mongo.Collection
	AnimalColl *mongo.Collection
}

func NewMongoPetRepository(db *mongo.Database) IPetRepository {
	return &mongoPetRepository{
		DB:         db,
		PetColl:    db.Collection("pet"),
		AnimalColl: db.Collection("animal"),
	}
}

func (repo *mongoPetRepository) GetPetInfo(petID string) (*domain.ApiPetInfo, error) {
	mongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return nil, BAD_PET_ID
	}

	dbInfo := new(domain.DBPetInfo)
	err = repo.PetColl.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(dbInfo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, NOT_FOUND
	} else if err != nil {
		return nil, err
	}

	return dbInfo.ToApi(), nil
}

func (repo *mongoPetRepository) GetAvatarBytes(petID string) ([]byte, error) {
	mongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return nil, BAD_PET_ID
	}

	var avaUrl struct {
		AvatarBytes []byte `bson:"avatar_url"`
	}

	opt := options.FindOne().SetProjection(bson.M{"avatar_url": 1, "_id": 0})
	err = repo.PetColl.FindOne(context.TODO(), bson.M{"_id": mongoID}, opt).Decode(&avaUrl)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []byte{}, nil
	} else if err != nil {
		return nil, err
	}

	return avaUrl.AvatarBytes, nil
}

func (repo *mongoPetRepository) IncrementAnimal(typeOfAnimal string, serviceID string) error {
	objID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		return err
	}

	normalizedType := strings.TrimSpace(strings.ToLower(typeOfAnimal))

	filter := bson.M{"type_of_animal": normalizedType}

	servDBRef := bson.M{
		"$ref": "service",
		"$id":  objID,
	}

	update := bson.D{
		{"$setOnInsert", bson.M{
			"type_of_animal": normalizedType,
			//"count":          0,
			//"services":       []bson.M{servDBRef},
		}},
		{"$inc", bson.M{"count": 1}},
		{"$push", bson.M{"services": servDBRef}},
		/*"$setOnInsert": bson.M{
			"type_of_animal": normalizedType,
			"count":          0,
			//"services":       []bson.M{servDBRef},
		},*/
	}

	opts := options.Update().SetUpsert(true)
	_, err = repo.AnimalColl.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (repo *mongoPetRepository) DecrementAnimal(typeOfAnimal string, serviceID string) error {
	objID, err := bson.ObjectIDFromHex(serviceID)
	if err != nil {
		return err
	}

	normalizedType := strings.TrimSpace(strings.ToLower(typeOfAnimal))

	filter := bson.M{"type_of_animal": normalizedType}

	servDBRef := bson.M{
		"$ref": "service",
		"$id":  objID,
	}

	update := bson.M{
		"$inc":  bson.M{"count": -1},
		"$pull": bson.M{"services": servDBRef},
	}

	var updatedDoc bson.M
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err = repo.AnimalColl.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedDoc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return NOT_FOUND
	} else if err != nil {
		return err
	}

	if count, ok := updatedDoc["count"].(int32); ok && count <= 0 {
		_, err := repo.AnimalColl.DeleteOne(context.TODO(), filter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *mongoPetRepository) GetTopAnimals(top int64) ([]string, error) {
	findOptions := options.Find().
		SetSort(bson.D{{"count", -1}}). // Sort by 'count' in descending order
		SetLimit(top).
		SetProjection(bson.M{"type_of_animal": 1, "_id": 0})

	cursor, err := repo.AnimalColl.Find(context.TODO(), bson.D{}, findOptions)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	stringResults := []string{}
	for _, res := range results {
		stringResults = append(stringResults, res["type_of_animal"].(string))
	}

	return stringResults, nil
}
