package mongoTLC

import (
	"context"
	"errors"
	"mainService/internal/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type IPetRepository interface {
	GetPetInfo(petID string) (*domain.ApiPetInfo, error)
	GetAvatarPath(petID string) (string, error)
}

type mongoPetRepository struct {
	DB   *mongo.Database
	Coll *mongo.Collection
}

func NewMongoPetRepository(db *mongo.Database) IPetRepository {
	return &mongoPetRepository{
		DB:   db,
		Coll: db.Collection("pet"),
	}
}

func (repo *mongoPetRepository) GetPetInfo(petID string) (*domain.ApiPetInfo, error) {
	mongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return nil, BAD_PET_ID
	}

	dbInfo := new(domain.DBPetInfo)
	err = repo.Coll.FindOne(context.TODO(), bson.M{"_id": mongoID}).Decode(dbInfo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, NOT_FOUND
	} else if err != nil {
		return nil, err
	}

	apiInfo := dbInfo.ToApi()
	if err != nil {
		return nil, err
	}

	return apiInfo, nil
}

func (repo *mongoPetRepository) GetAvatarPath(petID string) (string, error) {
	mongoID, err := bson.ObjectIDFromHex(petID)
	if err != nil {
		return "", BAD_PET_ID
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
