package domain

import (
	"encoding/base64"
	"mainService/pkg/serverErrors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Role string

const (
	Provider Role = "slave"
	Customer Role = "master"
)

func IsRole(str Role) bool {
	return str == Provider || str == Customer
}

type ApiService struct {
	ServiceID   string   `json:"service_id,omitempty"`
	Type        Role     `json:"role,omitempty"`
	UserID      string   `json:"user_id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Price       int32    `json:"price"`
	Description string   `json:"description,omitempty"`
	UserImage   string   `json:"user_image"`
	PetIDs      []string `json:"pet_ids"`
}

type DBService struct {
	ServiceID   bson.ObjectID `bson:"_id,omitempty"`
	Type        Role          `bson:"role,omitempty"`
	UserID      bson.M        `bson:"owner,omitempty"`
	Title       string        `bson:"title"`
	Price       int32         `bson:"price,omitempty"`
	Description string        `bson:"description,omitempty"`
	UserImage   []byte        `bson:"user_image,omitempty"`
	PetIDs      []bson.M      `bson:"pets,omitempty"`
}

func (api *ApiService) ToDB() (*DBService, error) {
	dbServ := &DBService{
		Type:        api.Type,
		Title:       api.Title,
		Description: api.Description,
		Price:       api.Price,
	}

	if api.ServiceID != "" {
		serviceID, err := bson.ObjectIDFromHex(api.ServiceID)
		if err != nil {
			return nil, err
		}

		dbServ.ServiceID = serviceID
	}

	if api.UserID != "" {
		mongoUserID, err := bson.ObjectIDFromHex(api.UserID)
		if err != nil {
			return nil, err
		}

		userDBRef := bson.M{
			"$ref": "user",
			"$id":  mongoUserID,
		}

		dbServ.UserID = userDBRef
	}

	if len(api.PetIDs) > 0 {
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
	}

	if api.UserImage != "" {
		byteImage, err := base64.StdEncoding.DecodeString(api.UserImage)
		if err != nil {
			return nil, err
		}

		dbServ.UserImage = byteImage
	}

	return dbServ, nil
}

func (db *DBService) ToApi() (*ApiService, error) {
	apiServ := &ApiService{
		ServiceID:   db.ServiceID.Hex(),
		Type:        db.Type,
		Title:       db.Title,
		Description: db.Description,
		Price:       db.Price,
	}

	if db.UserID != nil {
		userID, ok := db.UserID["$id"].(bson.ObjectID)
		if !ok {
			return nil, serverErrors.CAST_ERROR
		}

		apiServ.UserID = userID.Hex()
	}

	if len(db.PetIDs) > 0 {
		apiPetIDs := make([]string, len(db.PetIDs))
		for i, pet := range db.PetIDs {
			petID, ok := pet["$id"].(bson.ObjectID)
			if !ok {
				return nil, serverErrors.CAST_ERROR
			}
			apiPetIDs[i] = petID.Hex()
		}

		apiServ.PetIDs = apiPetIDs
	}

	if len(db.UserImage) != 0 {
		apiServ.UserImage = base64.StdEncoding.EncodeToString(db.UserImage)
	}

	return apiServ, nil
}

type DBServiceSerachResult struct {
	ServiceID   bson.ObjectID `bson:"_id,omitempty"`
	Type        Role          `bson:"role,omitempty"`
	UserID      bson.M        `bson:"owner,omitempty"`
	Title       string        `bson:"title"`
	Price       int32         `bson:"price,omitempty"`
	Description string        `bson:"description,omitempty"`
	Score       float64       `bson:"score,omitempty"`
	UserImage   []byte        `bson:"user_image,omitempty"`
	PetIDs      []bson.M      `bson:"pets,omitempty"`
}

func (db *DBServiceSerachResult) ToApiService() (*ApiService, error) {
	dbService := &DBService{
		ServiceID:   db.ServiceID,
		Type:        db.Type,
		UserID:      db.UserID,
		Title:       db.Title,
		Description: db.Description,
		Price:       db.Price,
		UserImage:   db.UserImage,
		PetIDs:      db.PetIDs,
	}

	return dbService.ToApi()
}

type ServiceFilter struct {
	MinPrice int32    `json:"min_price,omitempty"`
	MaxPrice int32    `json:"max_price,omitempty"`
	Animals  []string `json:"animals,omitempty"`
}
