package domain

import (
	"mainService/pkg/serverErrors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApiUserInfo struct {
	UserID        string   `json:"user_id,omitempty"`
	Login         string   `json:"login,omitempty"`
	Password      string   `json:"password,omitempty"`
	Username      string   `json:"username,omitempty"`
	Contacts      string   `json:"contacts,omitempty"`
	UserImage     string   `json:"user_image_string,omitempty"`
	UserBackImage string   `json:"background_image_string,omitempty"`
	PetIDs        []string `json:"pet_ids,omitempty"`
}

type DBUserInfo struct {
	UserID         bson.ObjectID `bson:"_id,omitempty"`
	Login          string        `bson:"login,omitempty"`
	HashedPassword []byte        `bson:"hashed_password,omitempty"`
	Salt           []byte        `bson:"salt"`
	Username       string        `bson:"name,omitempty"`
	Contacts       string        `bson:"contact,omitempty"`
	UserImage      string        `bson:"avatar_url,omitempty"`
	UserBackImage  string        `bson:"background_url,omitempty"`
	PetIDs         []bson.M      `bson:"pets,omitempty"`
}

func (apiInfo *ApiUserInfo) ToDB() (*DBUserInfo, error) {
	dbInfo := &DBUserInfo{
		Login:         apiInfo.Login,
		Username:      apiInfo.Username,
		Contacts:      apiInfo.Contacts,
		UserImage:     apiInfo.UserImage,
		UserBackImage: apiInfo.UserBackImage,
	}

	dbID, err := bson.ObjectIDFromHex(apiInfo.UserID)
	if err != nil {
		return nil, err
	}

	dbInfo.UserID = dbID

	dbIDs := make([]bson.M, len(apiInfo.PetIDs))
	for i, pet := range apiInfo.PetIDs {
		mongoID, err := bson.ObjectIDFromHex(pet)
		if err != nil {
			return nil, err
		}

		petDBRef := bson.M{
			"$ref": "pet",
			"$id":  mongoID,
		}

		dbIDs[i] = petDBRef
	}

	return dbInfo, nil
}

func (dbInfo *DBUserInfo) ToApi() (*ApiUserInfo, error) {
	apiInfo := &ApiUserInfo{
		UserID:        dbInfo.UserID.Hex(),
		Login:         dbInfo.Login,
		Username:      dbInfo.Username,
		Contacts:      dbInfo.Contacts,
		UserImage:     dbInfo.UserImage,
		UserBackImage: dbInfo.UserBackImage,
	}

	if len(dbInfo.PetIDs) > 0 {
		strPetIDs := make([]string, len(dbInfo.PetIDs))
		for i, pet := range dbInfo.PetIDs {
			petID, ok := pet["$id"].(bson.ObjectID)
			if !ok {
				return nil, serverErrors.CAST_ERROR
			}

			strPetIDs[i] = petID.Hex()
		}

		apiInfo.PetIDs = strPetIDs
	}

	return apiInfo, nil
}
