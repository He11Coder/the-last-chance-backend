package domain

import (
	"encoding/base64"
	"mainService/pkg/authUtils"
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
	Salt           []byte        `bson:"salt,omitempty"`
	Username       string        `bson:"name,omitempty"`
	Contacts       string        `bson:"contact,omitempty"`
	UserImage      []byte        `bson:"avatar_url,omitempty"`
	UserBackImage  []byte        `bson:"background_url,omitempty"`
	PetIDs         []bson.M      `bson:"pets,omitempty"`
}

func (apiInfo *ApiUserInfo) ToDB() (*DBUserInfo, error) {
	dbInfo := &DBUserInfo{
		Login:    apiInfo.Login,
		Username: apiInfo.Username,
		Contacts: apiInfo.Contacts,
	}

	if apiInfo.UserID != "" {
		dbID, err := bson.ObjectIDFromHex(apiInfo.UserID)
		if err != nil {
			return nil, err
		}

		dbInfo.UserID = dbID
	}

	if len(apiInfo.PetIDs) > 0 {
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

		dbInfo.PetIDs = dbIDs
	}

	if apiInfo.UserImage != "" {
		imageBytes, err := base64.StdEncoding.DecodeString(apiInfo.UserImage)
		if err != nil {
			return nil, err
		}

		dbInfo.UserImage = imageBytes
	}

	if apiInfo.UserBackImage != "" {
		backImageBytes, err := base64.StdEncoding.DecodeString(apiInfo.UserBackImage)
		if err != nil {
			return nil, err
		}

		dbInfo.UserBackImage = backImageBytes
	}

	return dbInfo, nil
}

func (dbInfo *DBUserInfo) ToApi() (*ApiUserInfo, error) {
	apiInfo := &ApiUserInfo{
		UserID:   dbInfo.UserID.Hex(),
		Login:    dbInfo.Login,
		Username: dbInfo.Username,
		Contacts: dbInfo.Contacts,
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

	if len(dbInfo.UserImage) != 0 {
		apiInfo.UserImage = base64.StdEncoding.EncodeToString(dbInfo.UserImage)
	}

	if len(dbInfo.UserBackImage) != 0 {
		apiInfo.UserBackImage = base64.StdEncoding.EncodeToString(dbInfo.UserBackImage)
	}

	return apiInfo, nil
}

type ApiUserUpdate struct {
	Login         string `json:"login,omitempty"`
	OldPassword   string `json:"old_password,omitempty"`
	NewPassword   string `json:"new_password,omitempty"`
	Username      string `json:"username,omitempty"`
	Contacts      string `json:"contacts,omitempty"`
	UserImage     string `json:"user_image_string,omitempty"`
	UserBackImage string `json:"background_image_string,omitempty"`
}

type DBUserUpdate struct {
	Login          string `bson:"login,omitempty"`
	HashedPassword []byte `bson:"hashed_password,omitempty"`
	Salt           []byte `bson:"salt,omitempty"`
	Username       string `bson:"name,omitempty"`
	Contacts       string `bson:"contact,omitempty"`
	UserImage      []byte `bson:"avatar_url,omitempty"`
	UserBackImage  []byte `bson:"background_url,omitempty"`
}

func (api *ApiUserUpdate) ToDB() (*DBUserUpdate, error) {
	db := &DBUserUpdate{
		Login:    api.Login,
		Username: api.Username,
		Contacts: api.Contacts,
	}

	if api.NewPassword != "" {
		newHashedPass, newSalt, err := authUtils.GenerateHash(api.NewPassword)
		if err != nil {
			return nil, err
		}

		db.HashedPassword = newHashedPass
		db.Salt = newSalt
	}

	if api.UserImage != "" {
		imageBytes, err := base64.StdEncoding.DecodeString(api.UserImage)
		if err != nil {
			return nil, err
		}

		db.UserImage = imageBytes
	}

	if api.UserBackImage != "" {
		backImageBytes, err := base64.StdEncoding.DecodeString(api.UserBackImage)
		if err != nil {
			return nil, err
		}

		db.UserBackImage = backImageBytes
	}

	return db, nil
}
