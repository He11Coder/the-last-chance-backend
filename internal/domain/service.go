package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Role string

const (
	Provider Role = "slave"
	Customer Role = "master"
)

type ApiService struct {
	Type        Role     `json:"role"`
	UserID      string   `json:"user_id,omitempty"`
	Title       string   `json:"title"`
	Description string   `json:"description,omitempty"`
	UserImage   string   `json:"user_image,omitempty"`
	PetIDs      []string `json:"pet_ids,omitempty"`
}

type DBService struct {
	Type        Role            `bson:"role"`
	UserID      bson.ObjectID   `bson:"_id,omitempty"`
	Title       string          `bson:"title"`
	Description string          `bson:"description,omitempty"`
	UserImage   string          `bson:"user_image,omitempty"`
	PetIDs      []bson.ObjectID `bson:"pet_ids,omitempty"`
}

func (api *ApiService) ToDB() (*DBService, error) {
	dbServ := &DBService{
		Type:        api.Type,
		Title:       api.Title,
		Description: api.Description,
		UserImage:   api.UserImage,
	}

	dbUserID, err := bson.ObjectIDFromHex(api.UserID)
	if err != nil {
		return nil, err
	}

	dbPetIDs := make([]bson.ObjectID, len(api.PetIDs))
	for i, apiID := range api.PetIDs {
		dbID, err := bson.ObjectIDFromHex(apiID)
		if err != nil {
			return nil, err
		}

		dbPetIDs[i] = dbID
	}

	dbServ.UserID = dbUserID
	dbServ.PetIDs = dbPetIDs

	return dbServ, nil
}

func (db *DBService) ToApi() *ApiService {
	apiServ := &ApiService{
		Type:        db.Type,
		UserID:      db.UserID.Hex(),
		Title:       db.Title,
		Description: db.Description,
		UserImage:   db.UserImage,
	}

	apiPetIDs := make([]string, len(db.PetIDs))
	for i, apiID := range db.PetIDs {
		apiPetIDs[i] = apiID.Hex()
	}

	apiServ.PetIDs = apiPetIDs

	return apiServ
}
