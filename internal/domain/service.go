package domain

import (
	"mainService/pkg/serverErrors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	Provider Role = "slave"
	Customer Role = "master"
)

type ApiService struct {
	ServiceID   string   `json:"service_id"`
	Type        Role     `json:"role"`
	UserID      string   `json:"user_id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	UserImage   string   `json:"user_image,omitempty"`
	PetIDs      []string `json:"pet_ids,omitempty"`
}

type DBService struct {
	ServiceID   bson.ObjectID `bson:"_id,omitempty"`
	Type        Role          `bson:"role,omitempty"`
	UserID      bson.M        `bson:"owner,omitempty"`
	Title       string        `bson:"title"`
	Description string        `bson:"description,omitempty"`
	UserImage   string        `bson:"user_image,omitempty"`
	PetIDs      []bson.M      `bson:"pets,omitempty"`
}

func (api *ApiService) ToDB() (*DBService, error) {
	dbServ := &DBService{
		Type:        api.Type,
		Title:       api.Title,
		Description: api.Description,
		UserImage:   api.UserImage,
	}

	serviceID, err := bson.ObjectIDFromHex(api.ServiceID)
	if err != nil {
		return nil, err
	}

	dbServ.ServiceID = serviceID

	mongoUserID, err := bson.ObjectIDFromHex(api.UserID)
	if err != nil {
		return nil, err
	}

	userDBRef := bson.M{
		"$ref": "user",
		"$id":  mongoUserID,
	}

	dbServ.UserID = userDBRef

	dbPetIDs := make([]bson.M, len(api.PetIDs))
	for i, apiID := range api.PetIDs {
		dbID, err := bson.ObjectIDFromHex(apiID)
		if err != nil {
			return nil, err
		}

		petDBRef := bson.M{
			"$ref": "pet",
			"$id":  dbID,
		}

		dbPetIDs[i] = petDBRef
	}

	dbServ.PetIDs = dbPetIDs

	return dbServ, nil
}

func (db *DBService) ToApi() (*ApiService, error) {
	apiServ := &ApiService{
		ServiceID:   db.ServiceID.Hex(),
		Type:        db.Type,
		Title:       db.Title,
		Description: db.Description,
		UserImage:   db.UserImage,
	}

	userID, ok := db.UserID["$id"].(bson.ObjectID)
	if !ok {
		return nil, serverErrors.CAST_ERROR
	}

	apiServ.UserID = userID.Hex()

	apiPetIDs := make([]string, len(db.PetIDs))
	for i, pet := range db.PetIDs {
		petID, ok := pet["$id"].(bson.ObjectID)
		if !ok {
			return nil, serverErrors.CAST_ERROR
		}
		apiPetIDs[i] = petID.Hex()
	}

	apiServ.PetIDs = apiPetIDs

	return apiServ, nil
}
