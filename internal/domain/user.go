package domain

import (
	"mainService/pkg/serverErrors"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ApiUserInfo struct {
	UserID    string   `json:"user_id,omitempty"`
	Username  string   `json:"username,omitempty"`
	Contacts  string   `json:"contacts,omitempty"`
	UserImage string   `json:"user_image,omitempty"`
	PetIDs    []string `json:"pet_ids,omitempty"`
}

type DBUserInfo struct {
	UserID    bson.ObjectID `bson:"_id,omitempty"`
	Username  string        `bson:"name,omitempty"`
	Contacts  string        `bson:"contact,omitempty"`
	UserImage string        `bson:"avatar_url,omitempty"`
	PetIDs    []bson.M      `bson:"pets,omitempty"`
}

func (apiInfo *ApiUserInfo) ToDB() (*DBUserInfo, error) {
	dbInfo := &DBUserInfo{
		Username:  apiInfo.Username,
		Contacts:  apiInfo.Contacts,
		UserImage: apiInfo.UserImage,
	}

	dbID, err := bson.ObjectIDFromHex(apiInfo.UserID)
	if err != nil {
		return nil, err
	}

	dbInfo.UserID = dbID

	return dbInfo, nil
}

func (dbInfo *DBUserInfo) ToApi() (*ApiUserInfo, error) {
	apiInfo := &ApiUserInfo{
		UserID:    dbInfo.UserID.Hex(),
		Username:  dbInfo.Username,
		Contacts:  dbInfo.Contacts,
		UserImage: dbInfo.UserImage,
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

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
